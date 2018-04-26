package food_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	dataNormalizer "github.com/tidepool-org/platform/data/normalizer"
	"github.com/tidepool-org/platform/data/types/food"
	testDataTypes "github.com/tidepool-org/platform/data/types/test"
	testErrors "github.com/tidepool-org/platform/errors/test"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/structure"
	structureValidator "github.com/tidepool-org/platform/structure/validator"
	"github.com/tidepool-org/platform/test"
)

func NewIngredient(ingredientArrayDepth int) *food.Ingredient {
	datum := food.NewIngredient()
	datum.Amount = NewAmount()
	datum.Brand = pointer.String(test.NewText(1, 100))
	datum.Code = pointer.String(test.NewText(1, 100))
	datum.Ingredients = NewIngredientArray(ingredientArrayDepth)
	datum.Name = pointer.String(test.NewText(1, 100))
	datum.Nutrition = NewNutrition()
	return datum
}

func CloneIngredient(datum *food.Ingredient) *food.Ingredient {
	if datum == nil {
		return nil
	}
	clone := food.NewIngredient()
	clone.Amount = CloneAmount(datum.Amount)
	clone.Brand = test.CloneString(datum.Brand)
	clone.Code = test.CloneString(datum.Code)
	clone.Ingredients = CloneIngredientArray(datum.Ingredients)
	clone.Name = test.CloneString(datum.Name)
	clone.Nutrition = CloneNutrition(datum.Nutrition)
	return clone
}

func NewIngredientArray(ingredientArrayDepth int) *food.IngredientArray {
	var datum *food.IngredientArray
	if ingredientArrayDepth--; ingredientArrayDepth > 0 {
		datum = food.NewIngredientArray()
		for count := 0; count < test.RandomIntFromRange(1, 3); count++ {
			*datum = append(*datum, NewIngredient(ingredientArrayDepth))
		}
	}
	return datum
}

func CloneIngredientArray(datumArray *food.IngredientArray) *food.IngredientArray {
	if datumArray == nil {
		return nil
	}
	clone := food.NewIngredientArray()
	for _, datum := range *datumArray {
		*clone = append(*clone, CloneIngredient(datum))
	}
	return clone
}

var _ = Describe("Ingredient", func() {
	It("IngredientBrandLengthMaximum is expected", func() {
		Expect(food.IngredientBrandLengthMaximum).To(Equal(100))
	})

	It("IngredientCodeLengthMaximum is expected", func() {
		Expect(food.IngredientCodeLengthMaximum).To(Equal(100))
	})

	It("IngredientNameLengthMaximum is expected", func() {
		Expect(food.IngredientNameLengthMaximum).To(Equal(100))
	})

	Context("ParseIngredient", func() {
		// TODO
	})

	Context("NewIngredient", func() {
		It("is successful", func() {
			Expect(food.NewIngredient()).To(Equal(&food.Ingredient{}))
		})
	})

	Context("Ingredient", func() {
		Context("Parse", func() {
			// TODO
		})

		Context("Validate", func() {
			DescribeTable("validates the datum",
				func(mutator func(datum *food.Ingredient), expectedErrors ...error) {
					datum := NewIngredient(3)
					mutator(datum)
					testDataTypes.ValidateWithExpectedOrigins(datum, structure.Origins(), expectedErrors...)
				},
				Entry("succeeds",
					func(datum *food.Ingredient) {},
				),
				Entry("amount missing",
					func(datum *food.Ingredient) { datum.Amount = nil },
				),
				Entry("amount invalid",
					func(datum *food.Ingredient) { datum.Amount.Units = nil },
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/amount/units"),
				),
				Entry("amount valid",
					func(datum *food.Ingredient) { datum.Amount = NewAmount() },
				),
				Entry("brand missing",
					func(datum *food.Ingredient) { datum.Brand = nil },
				),
				Entry("brand empty",
					func(datum *food.Ingredient) { datum.Brand = pointer.String("") },
					testErrors.WithPointerSource(structureValidator.ErrorValueEmpty(), "/brand"),
				),
				Entry("brand length; in range (upper)",
					func(datum *food.Ingredient) { datum.Brand = pointer.String(test.NewText(100, 100)) },
				),
				Entry("brand length; out of range (upper)",
					func(datum *food.Ingredient) { datum.Brand = pointer.String(test.NewText(101, 101)) },
					testErrors.WithPointerSource(structureValidator.ErrorLengthNotLessThanOrEqualTo(101, 100), "/brand"),
				),
				Entry("code missing",
					func(datum *food.Ingredient) { datum.Code = nil },
				),
				Entry("code empty",
					func(datum *food.Ingredient) { datum.Code = pointer.String("") },
					testErrors.WithPointerSource(structureValidator.ErrorValueEmpty(), "/code"),
				),
				Entry("code length; in range (upper)",
					func(datum *food.Ingredient) { datum.Code = pointer.String(test.NewText(100, 100)) },
				),
				Entry("code length; out of range (upper)",
					func(datum *food.Ingredient) { datum.Code = pointer.String(test.NewText(101, 101)) },
					testErrors.WithPointerSource(structureValidator.ErrorLengthNotLessThanOrEqualTo(101, 100), "/code"),
				),
				Entry("ingredients missing",
					func(datum *food.Ingredient) { datum.Ingredients = nil },
				),
				Entry("ingredients invalid",
					func(datum *food.Ingredient) { (*datum.Ingredients)[0] = nil },
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/ingredients/0"),
				),
				Entry("ingredients valid",
					func(datum *food.Ingredient) { datum.Ingredients = NewIngredientArray(3) },
				),
				Entry("name missing",
					func(datum *food.Ingredient) { datum.Name = nil },
				),
				Entry("name empty",
					func(datum *food.Ingredient) { datum.Name = pointer.String("") },
					testErrors.WithPointerSource(structureValidator.ErrorValueEmpty(), "/name"),
				),
				Entry("name length; in range (upper)",
					func(datum *food.Ingredient) { datum.Name = pointer.String(test.NewText(100, 100)) },
				),
				Entry("name length; out of range (upper)",
					func(datum *food.Ingredient) { datum.Name = pointer.String(test.NewText(101, 101)) },
					testErrors.WithPointerSource(structureValidator.ErrorLengthNotLessThanOrEqualTo(101, 100), "/name"),
				),
				Entry("nutrition missing",
					func(datum *food.Ingredient) { datum.Nutrition = nil },
				),
				Entry("nutrition invalid",
					func(datum *food.Ingredient) { datum.Nutrition.Carbohydrate.Units = nil },
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/nutrition/carbohydrate/units"),
				),
				Entry("nutrition valid",
					func(datum *food.Ingredient) { datum.Nutrition = NewNutrition() },
				),
				Entry("multiple errors",
					func(datum *food.Ingredient) {
						datum.Amount.Units = nil
						datum.Brand = pointer.String("")
						datum.Code = pointer.String("")
						(*datum.Ingredients)[0] = nil
						datum.Name = pointer.String("")
						datum.Nutrition.Carbohydrate.Units = nil
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/amount/units"),
					testErrors.WithPointerSource(structureValidator.ErrorValueEmpty(), "/brand"),
					testErrors.WithPointerSource(structureValidator.ErrorValueEmpty(), "/code"),
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/ingredients/0"),
					testErrors.WithPointerSource(structureValidator.ErrorValueEmpty(), "/name"),
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/nutrition/carbohydrate/units"),
				),
			)
		})

		Context("Normalize", func() {
			DescribeTable("normalizes the datum",
				func(mutator func(datum *food.Ingredient)) {
					for _, origin := range structure.Origins() {
						datum := NewIngredient(3)
						mutator(datum)
						expectedDatum := CloneIngredient(datum)
						normalizer := dataNormalizer.New()
						Expect(normalizer).ToNot(BeNil())
						datum.Normalize(normalizer.WithOrigin(origin))
						Expect(normalizer.Error()).To(BeNil())
						Expect(normalizer.Data()).To(BeEmpty())
						Expect(datum).To(Equal(expectedDatum))
					}
				},
				Entry("does not modify the datum",
					func(datum *food.Ingredient) {},
				),
				Entry("does not modify the datum; amount missing",
					func(datum *food.Ingredient) { datum.Amount = nil },
				),
				Entry("does not modify the datum; brand missing",
					func(datum *food.Ingredient) { datum.Brand = nil },
				),
				Entry("does not modify the datum; code missing",
					func(datum *food.Ingredient) { datum.Code = nil },
				),
				Entry("does not modify the datum; ingredients missing",
					func(datum *food.Ingredient) { datum.Ingredients = nil },
				),
				Entry("does not modify the datum; name missing",
					func(datum *food.Ingredient) { datum.Name = nil },
				),
				Entry("does not modify the datum; nutrition missing",
					func(datum *food.Ingredient) { datum.Nutrition = nil },
				),
			)
		})
	})

	Context("ParseIngredientArray", func() {
		// TODO
	})

	Context("NewIngredientArray", func() {
		It("is successful", func() {
			Expect(food.NewIngredientArray()).To(Equal(&food.IngredientArray{}))
		})
	})

	Context("IngredientArray", func() {
		Context("Parse", func() {
			// TODO
		})

		Context("Validate", func() {
			DescribeTable("validates the datum",
				func(mutator func(datum *food.IngredientArray), expectedErrors ...error) {
					datum := food.NewIngredientArray()
					mutator(datum)
					testDataTypes.ValidateWithExpectedOrigins(datum, structure.Origins(), expectedErrors...)
				},
				Entry("succeeds",
					func(datum *food.IngredientArray) {},
					structureValidator.ErrorValueEmpty(),
				),
				Entry("empty",
					func(datum *food.IngredientArray) { *datum = *food.NewIngredientArray() },
					structureValidator.ErrorValueEmpty(),
				),
				Entry("nil",
					func(datum *food.IngredientArray) { *datum = append(*datum, nil) },
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/0"),
				),
				Entry("single invalid",
					func(datum *food.IngredientArray) {
						invalid := NewIngredient(3)
						invalid.Brand = pointer.String("")
						*datum = append(*datum, invalid)
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueEmpty(), "/0/brand"),
				),
				Entry("single valid",
					func(datum *food.IngredientArray) {
						*datum = append(*datum, NewIngredient(3))
					},
				),
				Entry("multiple invalid",
					func(datum *food.IngredientArray) {
						invalid := NewIngredient(3)
						invalid.Brand = pointer.String("")
						*datum = append(*datum, NewIngredient(3), invalid, NewIngredient(3))
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueEmpty(), "/1/brand"),
				),
				Entry("multiple valid",
					func(datum *food.IngredientArray) {
						*datum = append(*datum, NewIngredient(3), NewIngredient(3), NewIngredient(3))
					},
				),
				Entry("multiple; length in range (upper)",
					func(datum *food.IngredientArray) {
						for len(*datum) < 100 {
							*datum = append(*datum, NewIngredient(1))
						}
					},
				),
				Entry("multiple; length out of range (upper)",
					func(datum *food.IngredientArray) {
						for len(*datum) < 101 {
							*datum = append(*datum, NewIngredient(1))
						}
					},
					structureValidator.ErrorLengthNotLessThanOrEqualTo(101, 100),
				),
				Entry("multiple errors",
					func(datum *food.IngredientArray) {
						invalid := NewIngredient(3)
						invalid.Brand = pointer.String("")
						*datum = append(*datum, nil, invalid, NewIngredient(3))
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/0"),
					testErrors.WithPointerSource(structureValidator.ErrorValueEmpty(), "/1/brand"),
				),
			)
		})

		Context("Normalize", func() {
			DescribeTable("normalizes the datum",
				func(mutator func(datum *food.IngredientArray)) {
					for _, origin := range structure.Origins() {
						datum := NewIngredientArray(3)
						mutator(datum)
						expectedDatum := CloneIngredientArray(datum)
						normalizer := dataNormalizer.New()
						Expect(normalizer).ToNot(BeNil())
						datum.Normalize(normalizer.WithOrigin(origin))
						Expect(normalizer.Error()).To(BeNil())
						Expect(normalizer.Data()).To(BeEmpty())
						Expect(datum).To(Equal(expectedDatum))
					}
				},
				Entry("does not modify the datum",
					func(datum *food.IngredientArray) {},
				),
				Entry("does not modify the datum; amount missing",
					func(datum *food.IngredientArray) { (*datum)[0].Amount = nil },
				),
				Entry("does not modify the datum; brand missing",
					func(datum *food.IngredientArray) { (*datum)[0].Brand = nil },
				),
				Entry("does not modify the datum; ingredients missing",
					func(datum *food.IngredientArray) { (*datum)[0].Ingredients = nil },
				),
				Entry("does not modify the datum; name missing",
					func(datum *food.IngredientArray) { (*datum)[0].Name = nil },
				),
				Entry("does not modify the datum; nutrition missing",
					func(datum *food.IngredientArray) { (*datum)[0].Nutrition = nil },
				),
			)
		})
	})
})
