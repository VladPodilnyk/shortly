/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./cmd/shortly/public/*.{html,js}"],
    theme: {
      extend: {},
    },
    plugins: [require("@tailwindcss/forms"), require("@tailwindcss/typography")],
  };