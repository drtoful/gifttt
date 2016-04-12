Installation and usage
----------------------

See [gopkg.in/twik.v1](https://gopkg.in/twik.v1) for documentation and usage details.

Changes
-------

These are the changes to the original code:

- added convenience functions ''unless'' and ''when'' that have the same semantic as in Common Lisp. ''if'' now always needs three expressions to function.
- added comparision functions ''>'', ''>='', ''<'' and ''<='' that can only be applied to numbers (int and float) and always contain 2 epxressions. will return true or false.
- Scope is now an interface and the previous was renamed to DefaultScope
- Scope's can now be "stacked" via the Enclosure method
- added ''split'' functions for strings
- added ''nth'' and ''length'' functions for easy usage of list from ''split''


ToDo
----

- Update test cases with new functions
