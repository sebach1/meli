local combinations = import 'combinations.jsonnet';
local fooId = 1;
local barId = 5;

local mkOpts(base, altOpt, k, v) = {
  alt: base { [k]: altOpt[k] },
  zero: base { [k]: v },
};

{
  foo: {
    local base = self.none,
    local altOpt = $.bar.none,
    none: {
      id: fooId,
      available_quantity: fooId * 10,
      price: fooId * 100,
      attribute_combinations: combinations.foo,
      picture_ids: ['fooPicId'],
    },
    id: mkOpts(base, altOpt, 'id', 0),
    price: mkOpts(base, altOpt, 'price', 0),
    available_quantity: mkOpts(base, altOpt, 'available_quantity', null),
    attribute_combinations: mkOpts(base, altOpt, 'attribute_combinations', []),
    picture_ids: mkOpts(base, altOpt, 'picture_ids', []),
  },

  bar: {
    local base = self.none,
    local altOpt = $.foo.none,
    none: {
      id: barId,
      available_quantity: barId * 10,
      price: barId * 100,
      attribute_combinations: combinations.bar,
      picture_ids: ['barPicId'],
    },
    id: mkOpts(base, altOpt, 'id', 0),
    price: mkOpts(base, altOpt, 'price', 0),
    available_quantity: mkOpts(base, altOpt, 'available_quantity', null),
    attribute_combinations: mkOpts(base, altOpt, 'attribute_combinations', []),
    picture_ids: mkOpts(base, altOpt, 'picture_ids', []),
  },

  zero: {},
}
