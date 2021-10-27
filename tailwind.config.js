module.exports = {
  purge: {
    content: ["./embed/**/*.html", "./embed/**/*.js"],
    safelist: ["emoji", "big-emoji", "font-bold", "bg-mentionBlurple", "text-blurple", "rounded", "p-1", "text-sm"],
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

      hoverGray: "#32353b",
      mentionBlurple: "#404675",
      mentionText: "#dee0fc",
    },
    extend: {
      fontFamily: {
        sans: ["Inter", "Open Sans", "Helvetica", "Arial", "sans-serif"],
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
};
