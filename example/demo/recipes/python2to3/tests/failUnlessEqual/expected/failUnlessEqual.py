def test(self):
    r = self.parse("if 1 fooze", 'r3')
    self.assertEqual(
        r.tree.toStringTree(),
        '(if 1 fooze)'
    )
