module.exports = {
  "transpileDependencies": [
    "vuetify"
  ],
  "devServer": {
    "proxy": {
      "/": {
        "target": "http://localhost:8888/"
      }
    }
  }
}