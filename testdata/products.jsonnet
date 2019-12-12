local variants = import 'variants.jsonnet';
local fooId = 10;
local barId = 50;

local mkOpts(base, altOpt, k, v) = {
  alt: base { [k]: altOpt[k] },
  zero: base { [k]: v },
};


{
  foo: {
    local base = self.none,
    local altOpt = $.bar.none,
    none: {
      id: '10',
      category_id: 'fooCategory',
      title: 'fooProduct',
      available_quantity: fooId * 10,
      variations: [variants.foo.none],
      initial_quantity: fooId * 100,
      sold_quantity: fooId * 1000,
      price: fooId * 10000,
    },
    id: mkOpts(base, altOpt, 'id', ''),
    category_id: mkOpts(base, altOpt, 'category_id', ''),
    variations: mkOpts(base, altOpt, 'variations', []),
    title: mkOpts(base, altOpt, 'title', ''),
    available_quantity: mkOpts(base, altOpt, 'available_quantity', null),
    price: mkOpts(base, altOpt, 'price', 0),
  },

  bar: {
    local base = self.none,
    local altOpt = $.foo.none,
    none: {
      id: '50',
      category_id: 'barCategory',
      variations: [variants.bar.none],
      title: 'barProduct',
      available_quantity: barId * 10,
      initial_quantity: barId * 100,
      sold_quantity: barId * 1000,
      price: barId * 10000,
    },
    id: mkOpts(base, altOpt, 'id', ''),
    category_id: mkOpts(base, altOpt, 'category_id', ''),
    variations: mkOpts(base, altOpt, 'variations', []),
    title: mkOpts(base, altOpt, 'title', ''),
    available_quantity: mkOpts(base, altOpt, 'available_quantity', null),
    price: mkOpts(base, altOpt, 'price', 0),
  },
  zero: {},
}
