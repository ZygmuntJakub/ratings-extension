function createRatingContainer() {
  const ratingContainer = document.createElement("div");
  ratingContainer.classList.add("mk-rating-container");
  return ratingContainer;
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

const getAllCards = (retry = 0) => {
  return new Promise(async (resolve) => {
    if (retry > 5) {
      return resolve([]);
    }
    const cards = document.getElementsByClassName("title-card");
    if (cards.length === 0) {
      await sleep(1000);
      return resolve(await getAllCards());
    }
    return resolve(cards);
  });
};

const updateRating = (cards) => {
  for (let i = 0; i < cards.length; i++) {
    const card = cards[i];
    if (card.querySelector(".mk-rating-container")) {
      continue;
    }
    console.log("adding rating to card", card);
    const ratingContainer = createRatingContainer();
    ratingContainer.innerHTML = "6/10 ⭐️";
    ratingContainer.style.position = "absolute";
    ratingContainer.style.top = "10px";
    ratingContainer.style.right = "10px";
    ratingContainer.style.backgroundColor = "rgba(0, 0, 0, 0.5)";
    ratingContainer.style.color = "white";
    ratingContainer.style.padding = "5px";
    ratingContainer.style.borderRadius = "5px";
    ratingContainer.style.fontSize = "12px";
    ratingContainer.style.fontWeight = "bold";

    card.appendChild(ratingContainer);
  }
};

getAllCards().then(updateRating);

const observer = new MutationObserver((mutations) => {
  getAllCards().then(updateRating);
});

observer.observe(document.body, {
  childList: true,
  subtree: true,
});
