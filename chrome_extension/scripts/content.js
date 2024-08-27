const DELAY = 500;

function createRatingContainer() {
  const ratingContainer = document.createElement("div");
  ratingContainer.classList.add("mk-rating-container");
  return ratingContainer;
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function getAllCards(retry = 0) {
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
}

async function updateRating({ cards, ratings }) {
  for (let i = 0; i < cards.length; i++) {
    const card = cards[i];
    if (card.querySelector(".mk-rating-container")) {
      continue;
    }
    const { value } = ratings.get(card);
    const ratingContainer = createRatingContainer();
    ratingContainer.innerHTML = `${value}/10 ⭐️`;
    ratingContainer.style.position = "absolute";
    ratingContainer.style.top = "10px";
    ratingContainer.style.right = "10px";
    ratingContainer.style.backgroundColor = "rgba(0, 0, 0, 0.5)";
    ratingContainer.style.color = "white";
    ratingContainer.style.padding = "5px";
    ratingContainer.style.borderRadius = "5px";
    ratingContainer.style.fontSize = "12px";
    ratingContainer.style.fontWeight = "bold";
    ratingContainer.style.transition = "opacity 1s";
    ratingContainer.style.opacity = 0;
    setTimeout(() => {
      ratingContainer.style.opacity = 1;
    }, DELAY);

    card.appendChild(ratingContainer);
  }
}

const data = new Map();

async function collectDataFromCards(cards) {
  for (let i = 0; i < cards.length; i++) {
    const card = cards[i];
    if (data.has(card)) continue;
    const title = card.getElementsByClassName("fallback-text")[0].innerText;
    const id = card
      .getElementsByTagName("a")[0]
      .href.split("/")[4]
      .split("?")[0];
    data.set(card, { title, id });
  }

  return { cards, data };
}

const ratings = new Map();
const queue = new Set();
let timeoutId = null;

async function fetchRatings({ cards, data }) {
  return new Promise((resolve) => {
    // filter already fetched ratings
    for (let i = 0; i < cards.length; i++) {
      const card = cards[i];
      if (ratings.has(card)) {
        continue;
      }
      queue.add(card);

      // clear pending timeout
      if (timeoutId) clearTimeout(timeoutId);

      // create new timeout window
      timeoutId = setTimeout(() => {
        for (const c of queue) {
          ratings.set(c, { value: Math.round(Math.random() * 10) });
          queue.delete(c);
        }
        resolve({ cards, data, ratings });
      }, DELAY);
    }
  });
}

const observer = new MutationObserver(() => {
  getAllCards()
    .then(collectDataFromCards)
    .then(fetchRatings)
    .then(updateRating);
});

observer.observe(document.body, {
  childList: true,
  subtree: true,
});
