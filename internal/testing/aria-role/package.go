// Package ariarole defines ARIA roles, and determines roles for elements.
//
// # About ARIA Roles
//
// [ARIA roles] provide semantic meaning to HTML elements, bringing a semanetic
// structure to an HTML document.
//
// This is especially important for users of screen readers, that can provide
// more meaningful information when reading the page.
//
// Many HTML elements have inherent roles, such as a <button> has the role
// "button", but you can assign any role to any element with the "role"
// attribute. e.g. a <select> is the default "combobox" role, you can create
// your own using <div role="combobox">.
//
// Landmark roles can also dramatically improve the browsing experience.
// Landmarks can indicate the major regions of a page, and allows the user to
// quickly skip irrelevant headers. Example elements are <header>, <nav>,
// <footer>, and <main>.
//
// When an HTML element has an inherent role, it's preferable to use the native
// element, as it brings default behaviour, e.g., a <button> element has the
// can receive tab focus, and call the click handler on keyboard input. A
// non-standard "button", like <div role="button"> need custom code to behave as
// a button.
//
// Some reasons for not using standard elements are:
//   - Customize the look & feel. This shouldn't be necessary with a <button> as
//     they are just inline elements with a default CSS applied. A <select>
//     doesn't allow customizing the drop-down list.
//   - Customize behaviour, such as filtering or searching in a drop-down list,
//     where options are fetched from a server based on user input.
//
// Documentation of different roles typically document the type of behaviour
// they are expected to have, e.g., using up/down arrows change selected values.
//
// [ARIA roles]: https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Roles
package ariarole
