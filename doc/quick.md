# Quick Start

Start by creating a new file "my_first_rule.rule" with the following content

    (when (== time:second 5) (log "Hello World!"))

Now start giftt by invoking the binary. You will notice that every time your clock hits 5 seconds past the minute it will print out "Hello World!". You will also notice that it executes this rule every second. This is because the symbol "time:second" has changed its value. Every time a symbol changes its value, all rules that use this symbol will be evaluated. If a symbol does not change nothing will happen.

Let's create a second rule called "my_second_rule.rule" with the following content

    (when (== foo "bar") (log "Hello foobar!"))

Again you can start gifttt by invoking its binary. But unfortunately your new rule is never executed. You can change a symbol's value by calling the API endpoint and set a new value. So for example you can invoke the rule with:

     curl --data '{"value":"bar"}' http://localhost:4200/v/foo

You will see that your second rule has now been evaluated. If you call the API endpoint again with the same command, nothing will happen though. This is because there was no change in the symbol. Remember rules are only evaluated if a symbol changes.

Instead of setting a symbol's value using the API endpoint you can also set the value within another rule:

     (set foo "bar")

This will have the same effect as calling the endpoint. This allows you to create rule chains that have dependencies on other rules and symbol values.

There are a few symbols that are pre-defined that will help you in creating rules:

     time:second
     time:minute
     time:hour
     date:day
     date:month
     date:year

Values of symbols are persisted after they have been set. So you can safely stop gifttt and restart it afterwards to retain its internal state.
