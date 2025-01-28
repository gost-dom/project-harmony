/** @type {import('tailwindcss').Config} */
const colors = require("tailwindcss/colors");

module.exports = {
  content: ["./internal/server/views/*.templ"],
  theme: {
    extend: {},
    colors: {
      ...colors,
      primary: colors.slate,
      secondary: colors.orange,
      ctabase: colors.indigo,
      cta: colors.indigo[600],
    },
  },
  plugins: [],
};
