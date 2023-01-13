package main;
// point class comment
record Point<T>(int x, int y) {
    private static final int staticField = 15;

    Point(int x, int y, int z) {
        this(x, y);
        System.out.println(z);
    }

    public int x() {
        // comment
        int x = 12;
        return x;
    }

    public int x(int a) {
        // comment
        return x + a;
    }

    static {
        System.out.println("test");
    }
}
