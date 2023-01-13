class Test {
    Test() {
        //print(!(! true))
        print(/*comment1*/!/*comment2*/(!/*comment3*/true/*comment4*/)/*comment5*/);
        //print(!(! true && false))
        print(!(!true && false));
        //print(!(! true || false))
        print(!(!true || false));
        //print(! ! true)
        print(!!true);
        //print(! ! false)
        print(!!false);
        //print(! ! ! true)
        print(!!! true);
        //print(! ! ! ! true)
        print(!!!!true);
        //print(!(((! ! ! true))))
        print(!(((!!!true))));

        //print(!(!! true || false))
        print(!(!!true || false));

        //if (!(! testWhatever(12))) return "Hello"
        if (!(!testWhatever(12))) return "Hello";
    }
}
