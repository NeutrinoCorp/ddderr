package ddderr

import (
	"strconv"
	"strings"
)

const (
	// error types groups
	domain         = "Domain"
	infrastructure = "Infrastructure"
	// specific error types
	notFound      = "NotFound"
	alreadyExists = "AlreadyExists"
	outOfRange    = "OutOfRange"
	invalidFormat = "InvalidFormat"
	required      = "Required"
	remoteCall    = "FailedRemoteCall"

	unknownDomain         = "UnknownDomain"
	unknownInfrastructure = "UnknownInfrastructure"
)

// Error contains specific mechanisms useful for further error mapping and other
// specific use cases
type Error struct {
	parent      error
	group       string
	kind        string
	property    string
	title       string
	description string
	statusName  string

	dynamicDescription bool
	dynamicStatus      bool
	limitA, limitB     int
	formats            []string
}

var _ error = Error{}

// Error returns the error description
func (e Error) Error() string {
	return e.Description()
}

// Kind retrieves the error type (e.g. NotFound, AlreadyExists)
func (e Error) Kind() string {
	return e.kind
}

// SetKind sets the specific error type (e.g. NotFound, AlreadyExists)
func (e Error) SetKind(kind string) Error {
	e.kind = kind
	return e
}

// Property returns the resource or field which contains the error
func (e Error) Property() string {
	return e.property
}

// SetProperty sets the field or resource for an error
func (e Error) SetProperty(property string) Error {
	e.property = property
	e.dynamicDescription = true
	e.dynamicStatus = true
	return e
}

// Title retrieves a generic error message
func (e Error) Title() string {
	return e.title
}

// SetTitle sets a generic error message
func (e Error) SetTitle(title string) Error {
	e.title = title
	return e
}

// Description retrieves a specific and detailed error message
func (e Error) Description() string {
	if !e.dynamicDescription {
		return e.description
	}

	switch e.kind {
	case alreadyExists:
		return newAlreadyExistsDescription(e.property)
	case invalidFormat:
		if e.formats == nil {
			return e.description
		}
		return newInvalidFormatDescription(e.property, e.formats...)
	case remoteCall:
		return newRemoteCallDescription(e.property)
	case notFound:
		return newNotFoundDescription(e.property)
	case outOfRange:
		return newOutOfRangeDescription(e.property, e.limitA, e.limitB)
	case required:
		return newRequiredDescription(e.property)
	default:
		return e.description
	}
}

// SetDescription sets a specific and detailed error message
func (e Error) SetDescription(description string) Error {
	e.description = description
	e.dynamicDescription = false
	return e
}

// SetStatus sets a system-owned status for the Error
func (e Error) SetStatus(status string) Error {
	e.statusName = status
	e.dynamicStatus = false
	return e
}

// Status retrieves the status name of the Error owned by the system
func (e Error) Status() string {
	if !e.dynamicStatus {
		return e.statusName
	}
	return getSanitizedStatusName(e.property, e.kind)
}

// Parent returns the error parent
//
// Note: Might return nil if parent was not specified
func (e Error) Parent() error {
	return e.parent
}

// SetParent sets a parent error to the given DDD error
func (e Error) SetParent(err error) Error {
	e.parent = err
	return e
}

// IsDomain checks if the error belongs to Domain error group
func (e Error) IsDomain() bool {
	return e.group == domain
}

// IsInfrastructure checks if the error belongs to Infrastructure error group
func (e Error) IsInfrastructure() bool {
	return e.group == infrastructure
}

// IsRemoteCall checks if the error belongs to Failed Remote Call error types
func (e Error) IsRemoteCall() bool {
	return e.kind == remoteCall
}

// IsNotFound checks if the error belongs to Not Found error types
func (e Error) IsNotFound() bool {
	return e.kind == notFound
}

// IsAlreadyExists checks if the error belongs to Already Exists error types
func (e Error) IsAlreadyExists() bool {
	return e.kind == alreadyExists
}

// IsOutOfRange checks if the error belongs to Out of Range error types
func (e Error) IsOutOfRange() bool {
	return e.kind == outOfRange
}

// IsInvalidFormat checks if the error belongs to Invalid Format error types
func (e Error) IsInvalidFormat() bool {
	return e.kind == invalidFormat
}

// IsRequired checks if the error belongs to Required error types
func (e Error) IsRequired() bool {
	return e.kind == required
}

// NewDomain creates an Error for Domain generic use cases
func NewDomain(title, description string) Error {
	return Error{
		parent:      nil,
		group:       domain,
		kind:        unknownDomain,
		property:    "",
		title:       title,
		description: description,
	}
}

// NewInfrastructure creates an Error for Infrastructure generic use cases
func NewInfrastructure(title, description string) Error {
	return Error{
		parent:      nil,
		group:       infrastructure,
		kind:        unknownInfrastructure,
		property:    "",
		title:       title,
		description: description,
	}
}

// NewRemoteCall creates an Error for network remote calls failing scenarios
//
// (e.g. database connection failed, sync inter-service transaction failed over a networking problem)
func NewRemoteCall(externalResource string) Error {
	return Error{
		parent:      nil,
		group:       infrastructure,
		kind:        remoteCall,
		property:    externalResource,
		title:       "Remote call failed",
		description: newRemoteCallDescription(externalResource),
		statusName:  "FailedRemoteCall",
	}
}

func newRemoteCallDescription(resource string) string {
	desc := "Failed to call external resource"
	if resource != "" {
		desc = desc + " [" + resource + "]"
	}
	return desc
}

// NewNotFound creates an Error for Not Found use cases
//
// (description e.g. The resource foo was not found)
func NewNotFound(resource string) Error {
	return Error{
		parent:      nil,
		group:       domain,
		kind:        notFound,
		property:    resource,
		title:       "Resource not found",
		description: newNotFoundDescription(resource),
		statusName:  getSanitizedStatusName(resource, "NotFound"),
	}
}

func newNotFoundDescription(resource string) string {
	desc := "not found"
	if resource != "" {
		desc = "The resource " + resource + " was not found"
	}
	return desc
}

// NewAlreadyExists creates an Error for Already Exists use cases
//
// (description e.g. The resource foo was already created)
func NewAlreadyExists(resource string) Error {
	return Error{
		parent:      nil,
		group:       domain,
		kind:        alreadyExists,
		property:    resource,
		title:       "Resource already exists",
		description: newAlreadyExistsDescription(resource),
		statusName:  getSanitizedStatusName(resource, "AlreadyExists"),
	}
}

func newAlreadyExistsDescription(resource string) string {
	desc := "already exists"
	if resource != "" {
		desc = "The resource " + resource + " already exists"
	}
	return desc
}

// NewOutOfRange creates an Error for Out of Range use cases
//
// (description e.g. The property foo is out of range [A, B))
func NewOutOfRange(property string, a, b int) Error {
	return Error{
		parent:      nil,
		group:       domain,
		kind:        outOfRange,
		property:    property,
		title:       "Property is out of the specified range",
		description: newOutOfRangeDescription(property, a, b),
		statusName:  getSanitizedStatusName(property, "OutOfRange"),
		limitA:      a,
		limitB:      b,
	}
}

func newOutOfRangeDescription(property string, a, b int) string {
	desc := "out of range [" + strconv.Itoa(a) + "," + strconv.Itoa(b) + ")"
	if property != "" {
		desc = "The property " + property + " is " + desc
	}
	return desc
}

// NewInvalidFormat creates an Error for Invalid Format use cases
//
// (description e.g. The property foo has an invalid format, expected [x1, x2, xN])
func NewInvalidFormat(property string, formats ...string) Error {
	return Error{
		parent:      nil,
		group:       domain,
		kind:        invalidFormat,
		property:    property,
		title:       "Property is not a valid format",
		description: newInvalidFormatDescription(property, formats...),
		statusName:  getSanitizedStatusName(property, "InvalidFormat"),
		formats:     formats,
	}
}

func newInvalidFormatDescription(property string, formats ...string) string {
	desc := "invalid format, expected [" + strings.Join(formats, ",") + "]"
	if property != "" {
		desc = "The property " + property + " has an " + desc
	}
	return desc
}

// NewRequired creates an Error for Required use cases
//
// (description e.g. The property foo is required)
func NewRequired(property string) Error {
	return Error{
		parent:      nil,
		group:       domain,
		kind:        required,
		property:    property,
		title:       "Missing property",
		description: newRequiredDescription(property),
		statusName:  getSanitizedStatusName(property, "IsRequired"),
	}
}

func newRequiredDescription(property string) string {
	desc := "required"
	if property != "" {
		desc = "The property " + property + " is " + desc
	}
	return desc
}
