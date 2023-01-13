package main;

// point class comment
class Point<T> {
    private final int x;
    private final int y;

    private static final int staticField = 15;

    Point(int x, int y) {
        this.x = x;
        this.y = y;
    }


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

    public int y() {
        return y;
    }


    static {
        System.out.println("test");
    }
}
