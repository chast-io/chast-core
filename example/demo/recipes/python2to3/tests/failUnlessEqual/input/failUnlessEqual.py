def test(self):
    r = self.parse("if 1 fooze", 'r3')
    self.failUnlessEqual(
        r.tree.toStringTree(),
        '(if 1 fooze)'
    )
