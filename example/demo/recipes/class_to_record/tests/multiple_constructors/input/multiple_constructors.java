class Point {
    private final int x;
    private final int y;

    Point(int x, int y) {
        this.x = x;
        this.y = y;
    }

    Point() {
        this(0,0);
    }

    public int x() {
        return x;
    }

    public int y() {
        return y;
    }
}
