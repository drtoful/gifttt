# Rules

Rules are plain-text files ending with ".rule" and contain valid Lisp S-Expressions. These are statements in the form "(A B C)". Where A, B and C can themselves be S-Expressions (so you can nest expressions indefinetely). Altough you can write everything in one line, we suggest to use sensible line-breaking and formatting to make code more readable. Most rules will be hopefully short, so that the parenthesis are kept at a bare minimum.

The following gives a short summary, which predefined symbols (everything in Lisp is a symbol) there exists and how you can use them. The language is somewhat loosely typed. You can overwrite values with different typed ones at any time, so you are responsible yourself when creating rules, to keep track what type a symbol should have, otherwise you will run into run-time exceptions. The types are as following:

* string: any quoted text
* integer: any number that is not a float
* float: a number that contains a "."
* boolean: represented by the symbols **true** and **false**
* **nil**: the null value

## Reference

### log

    (log <text>)

Prints *text* into the application log. Always evaluates to **nil**.

### run

    (run <cmd> [<args>...])

Executes *cmd* with any number of arguments. Always evaluates to **nil**.

### when

    (when <condition> <action>)

If *condition* evaluates to **true** then the expression in *action* will be evaluated. Always evaluates to **nil**

### unless

    (unless <condition> <action>)

If *condition* evaluates to **false** then the epxression in *action* will be evaluated. Always evaluates to **nil**

### if

    (if <condition> <then_action> <else_action>)

If *condition* evaluates to **true** then the expression in *then_action* will be evaluated, *else_action* otherwise. Always evaluates to **nil**.

### do

    (do <expression> [<expression>...])

Evaluates all expressions in order. Always evaluates to **nil**.

### Comparators

    (== <a> <b>)
    (!= <a> <b>)

Compares expression *a* with expression *b*. Both need to be of same type. Always return **true** or **false**.

### Boolean Operators

    (and <a> [<b>...])
    (or <a> [<b>...])

Operates on boolean expressions. Always returns **true** or **false**

### Numeric Comparators

    (> <a> <b>)
    (>= <a> <b>)
    (< <a> <b>)
    (<= <a> <b>)

Compares two numeric expressions *a* and *b*. The types have to be either integer or float, however they can both be of different type. The comparision will always be conducted as float by comparing **a-b >= 0**. Always return **true** or **false**.

### Numeric Operators

    (+ <a> [<b>...])
    (* <a> [<b>...])
    (- <a> [<b>...])
    (/ <a> [<b>...])

Computes the new numeric value of the operation. The type of *a* and *b* have to be integer or float and can be mixed. If the type is different, the operation will always be done in float and the returned value will have type float. Otherwise the returned type will have the same type as the evaluated expressions.

### var

    (var <name> <value>)

Create a new symbol with name and value given. Note, that this will behave like a local variable. Setting this symbol to a different value will never trigger other rules (as the global symbol is never changed, and thus will not trigger the re-evaluation of all rules). This is handy to create local constants to keep the code more readable. Always evaluates to **nil**.

### set

    (set <name> <value>)

Set the value of a named symbol. If the symbol is global this will trigger re-evaluation of all rules, that contain this symbol. Always evaluates to **nil**.
