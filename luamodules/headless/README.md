# Headless

Headless chrome Lua port.

If you don't have a headless Chrome installed
it will be downloaded and installed automatically on the first launch.

## Types

- Headless
- Browser
- Page
- Element

### Headless

| Function | Return  |
|----------|---------|
| browser  | browser |

### Browser

| Function | Arguments                 | Return |
|----------|---------------------------|--------|
| page     | __optional__ string (url) | page   |

### Page

| Function             | Arguments      | Description                                                                                                                                                                                                                           | Return    |
|----------------------|----------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------|
| waitLoad             |                | waits for the `window.onload` event, it returns immediately if the event is already fired.                                                                                                                                            |           |
| element              | string         | retries until an element in the page that matches the CSS selector, then returns the matched element.                                                                                                                                 | element   |
| elementR             | string, string | retries until an element in the page that matches the css selector and it's text matches the jsRegex, then returns the matched element.                                                                                               | element   |
| elements             | string         | returns all elements that match the css selector                                                                                                                                                                                      | []element |
| elementByJS          | string         | returns the element from the return value of the js function. If sleeper is nil, no retry will be performed. By default, it will retry until the js function doesn't return null. To customize the retry logic, check the examples of | element   |
| waitElementsMoreThan | string, number | Wait until there are more than <num> <selector> elements.                                                                                                                                                                             |           |
| navigate             | string         | Navigate to the url. If the url is empty, "about:blank" will be used. It will return immediately after the server responds the http header.                                                                                           |           |
| has                  | string         | Has an element that matches the css selector                                                                                                                                                                                          | boolean   |
| eval                 | string         | Evaluate js function on the page. Will return a value if any is returned                                                                                                                                                              | string    |
| html                 | string         | HTML of the page                                                                                                                                                                                                                      | string    |

### Element

| Function  | Arguments | Description                                                                                                                                                                                                 | Return |
|-----------|-----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|
| input     | string    | Input focuses on the element and input text to it. Before the action, it will scroll to the element, wait until it's visible, enabled and writable. To empty the input you can use something like input("") |        |
| click     |           | Click will press then release the button just like a human. Before the action, it will try to scroll to the element, hover the mouse over it, wait until the it's interactable and enabled.                 |        |
| text      |           | Text that the element displays                                                                                                                                                                              | string |
| attribute | string    | Attribute of the element.                                                                                                                                                                                   | string |
| html      |           | HTML of the element                                                                                                                                                                                         | string |
| property  | string    | Property of the DOM object                                                                                                                                                                                  | string |
|           |           |                                                                                                                                                                                                             |        |

> [Property vs Attribute](https://stackoverflow.com/questions/6003819/what-is-the-difference-between-properties-and-attributes-in-html)

## Example

```lua
local headless = require("headless")
local browser = headless.browser()
local page = browser:page()
page:navigate("https://www.google.com")
page:waitLoad()
local element = page:element("input[name='q']")
element:input("lua")
local button = page:element("button[name='btnK']")
button:click()
page:waitLoad()
print(page:html())
```

> This example was generated by GitHub copilot, lol. I didn't test but looks legit.