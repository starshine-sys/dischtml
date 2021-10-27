module.exports = {
  purge: {
    content: ["./embed/**/*.html", "./embed/**/*.js"],
    safelist: ["emoji"],
  },
  darkMode: false, // or 'media' or 'class'
  theme: {
    colors: {
      blurple: "#5865F2",
      green: "#57F287",
      yellow: "#FEE75C",
      fuchsia: "#EB459E",
      red: "#ED4245",
      white: "#FFFFFF",
      black: "#000000",
      lightGray: "#36393F",
      darkGray: "#2F3136",
      lighterGray: "#bcbdbf",
      linkColour: "#00AFF4",

      hoverGray: "32353b",
    },
    extend: {
      fontFamily: {
        'sans': ['Inter', 'Open Sans', 'Helvetica', 'Arial', 'sans-serif']
      }
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
};
