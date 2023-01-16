import org.antlr.v4.runtime.CharStream
import org.antlr.v4.runtime.Lexer
import java.util.*

internal abstract class CSharpLexerBase protected constructor(input: CharStream?) : Lexer(input) {
    protected var interpolatedStringLevel = 0
    protected var interpolatedVerbatiums = Stack<Boolean>()
    protected var curlyLevels = Stack<Int>()
    protected var verbatium = false
    protected fun OnInterpolatedRegularStringStart() {
        interpolatedStringLevel++
        interpolatedVerbatiums.push(false)
        verbatium = false
    }

    protected fun OnInterpolatedVerbatiumStringStart() {
        interpolatedStringLevel++
        interpolatedVerbatiums.push(true)
        verbatium = true
    }

    protected fun OnOpenBrace() {
        if (interpolatedStringLevel > 0) {
            curlyLevels.push(curlyLevels.pop() + 1)
        }
    }

    protected fun OnCloseBrace() {
        if (interpolatedStringLevel > 0) {
            curlyLevels.push(curlyLevels.pop() - 1)
            if (curlyLevels.peek() == 0) {
                curlyLevels.pop()
                skip()
                popMode()
            }
        }
    }

    protected fun OnColon() {
        if (interpolatedStringLevel > 0) {
            var ind = 1
            var switchToFormatString = true
            while (_input.LA(ind).toChar() != '}') {
                if (_input.LA(ind) == ':'.code || _input.LA(ind) == ')'.code) {
                    switchToFormatString = false
                    break
                }
                ind++
            }
            if (switchToFormatString) {
                mode(CSharpLexer.INTERPOLATION_FORMAT)
            }
        }
    }

    protected fun OpenBraceInside() {
        curlyLevels.push(1)
    }

    protected fun OnDoubleQuoteInside() {
        interpolatedStringLevel--
        interpolatedVerbatiums.pop()
        verbatium = if (interpolatedVerbatiums.size > 0) interpolatedVerbatiums.peek() else false
    }

    protected fun OnCloseBraceInside() {
        curlyLevels.pop()
    }

    protected fun IsRegularCharInside(): Boolean {
        return !verbatium
    }

    protected fun IsVerbatiumDoubleQuoteInside(): Boolean {
        return verbatium
    }
}
