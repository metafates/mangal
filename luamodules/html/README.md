# HTML

Goquery lua port

## Types

- HTML
- Document
- Selection

### HTML

| Function | Arguments | Return   |
|----------|-----------|----------|
| parse    | string    | document |

### Document

| Function | Arguments | Return    |
|----------|-----------|-----------|
| find     | string    | selection |

### Selection

| Function | Arguments                | Description                                                                                                                                                                                                                                                                                               | Return    |
|----------|--------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------|
| find     | string                   | Find gets the descendants of each element in the current set of matched elements, filtered by a selector. It returns a new Selection object containing these matched elements.                                                                                                                            | selection |
| each     | function(num, selection) | Each iterates over a Selection object, executing a function for each matched element. It returns the current Selection object. The function f is called for each element in the selection with the index of the element in that selection starting at 0, and a Selection that contains only that element. |           |
| attr     | string                   | Attr gets the specified attribute's value for the first element in the Selection.                                                                                                                                                                                                                         | string    |
| first    |                          | First reduces the set of matched elements to the first in the set. It returns a new Selection object, and an empty Selection object if the the selection is empty.                                                                                                                                        | selection |
| parent   |                          | Parent gets the parent of each element in the Selection. It returns a new Selection object containing the matched elements.                                                                                                                                                                               | selection |
| text     |                          | Text gets the combined text contents of each element in the set of matched elements, including their descendants.                                                                                                                                                                                         | string    |
| html     |                          | Html gets the HTML contents of the first element in the set of matched elements. It includes text and comment nodes.                                                                                                                                                                                      | string    |
| hasClass | string                   | HasClass determines whether any of the matched elements are assigned the given class.                                                                                                                                                                                                                     | boolean   |
| is       | string                   | Is checks the current matched set of elements against a selector and returns true if at least one of these elements matches.                                                                                                                                                                              | boolean   |
| next     |                          | Next gets the immediately following sibling of each element in the Selection. It returns a new Selection object containing the matched elements.                                                                                                                                                          | selection |
| prev     |                          | Prev gets the immediately preceding sibling of each element in the Selection. It returns a new Selection object containing the matched elements.                                                                                                                                                          | selection |

## Example

```lua
local html = require("html")

local doc = html.parse("...")
doc:find(".classname > .anotherclass"):each(function(i, el)
    print(el:text())
end)
```