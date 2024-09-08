chrome.runtime.onMessage.addListener(function (message, _, senderResponse) {
  if (message.type === "ratings") {
    fetch(
      encodeURI(
        `https://master.bieda.it/api/v1/ratings?query=${message.query}`,
      ),
      {
        method: "GET",
      },
    )
      .then((res) => {
        return res.json();
      })
      .then((res) => {
        senderResponse(res);
      })
      .catch(console.log);
  }
  return true;
});
