package validator_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/ScarletTanager/validator"
)

var _ = Describe("Validator", func() {
	type OtherStructWithTags struct {
		Profession string `validator:"allowempty"`
	}

	type StructWithTags struct {
		First          string   `validator:"required"`
		Last           string   `validator:"required,allowempty"`
		Weight         int      `validator:"required,greaterthan,0"`
		ZipCode        string   `validator:"required,format(ddddd)"`
		Nicknames      []string `validator:"required"`
		EmbeddedStruct OtherStructWithTags
	}

	var (
		thing StructWithTags
	)

	BeforeEach(func() {
		thing = StructWithTags{
			First:   "foo",
			Last:    "bar",
			Weight:  250,
			ZipCode: "90210",
		}
	})

	Context("When no validations fail", func() {
		BeforeEach(func() {
			thing.Nicknames = []string{"foo", "frank", "george"}
		})

		It("Returns an empty list of errors", func() {
			Expect(Validate(thing)).To(HaveLen(0))
		})
	})

	Context("When attempting to validate a non-Struct Kind", func() {
		It("Returns an appropriate error", func() {
			Expect((Validate("foo"))[0]).To(MatchError("Incorrect Kind: string, must be a reflect.Struct"))
		})
	})

	Describe("String fields", func() {
		Context("When a required field is the empty string", func() {
			BeforeEach(func() {
				thing.First = ""
			})

			It("Returns a list containing the correct error", func() {
				Expect(Validate(thing)).To(ContainElement(&ValidationError{
					ErrorType: ValidationFailed,
					Message:   fmt.Sprintf(ErrorMessageValidationFailed, "First"),
				}))
			})
		})

		Context("When a required field with allowempty is set to the empty string", func() {
			BeforeEach(func() {
				thing.Last = ""
			})

			It("Accepts the empty string", func() {
				Expect(Validate(thing)).NotTo(ContainElement(&ValidationError{
					ErrorType: ValidationFailed,
					Message:   fmt.Sprintf(ErrorMessageValidationFailed, "Last"),
				}))
			})
		})

		Describe("Describe format specifiers", func() {
		})
	})

	Describe("Int fields", func() {
		Describe("Greater than", func() {
			Context("When the value is equal to the boundary param", func() {
				BeforeEach(func() {
					thing.Weight = 0
				})

				It("Returns a list containing the correct error", func() {
					Expect(Validate(thing)).To(ContainElement(&ValidationError{
						ErrorType: ValidationFailed,
						Message:   fmt.Sprintf(ErrorMessageValidationFailed, "Weight"),
					}))
				})
			})
		})
	})

	Describe("Slice fields", func() {

	})
})
