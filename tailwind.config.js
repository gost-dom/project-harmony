/** @type {import('tailwindcss').Config} */
const colors = require("tailwindcss/colors");

module.exports = {
  content: ["./internal/web/server/views/*.templ"],
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
