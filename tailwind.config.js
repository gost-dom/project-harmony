/** @type {import('tailwindcss').Config} */
const colors = require("tailwindcss/colors");

module.exports = {
  content: ["./views/*.templ"],
  theme: {
    extend: {},
    colors: {
      ...colors,
      primary: colors.slate,
      secondary: colors.orange,
      cta: colors.indigo[600],
    },
  },
  plugins: [],
};
