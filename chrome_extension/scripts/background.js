chrome.runtime.onMessage.addListener(function (message, _, senderResponse) {
  if (message.type === "ratings") {
    fetch(
      encodeURI(`http://localhost:8080/api/v1/ratings?query=${message.query}`),
      {
        method: "GET",
      },
    )
      .then((res) => {
        return res.json();
      })
      .then((res) => {
        senderResponse(res);
      });
  }
  return true;
});
