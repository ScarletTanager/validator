package validator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/ScarletTanager/validator"
)

var _ = Describe("Validator", func() {
	type StructWithTags struct {
		First   string `validator:"required"`
		Last    string `validator:"required,nonzero"`
		Weight  int    `validator:"required,greaterthan,0"`
		ZipCode string `validator:"required,format(ddddd)"`
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
				Expect(Validate(thing)).To(ContainElement(ValidationError{FieldName: "First"}))
			})
		})
	})

})
