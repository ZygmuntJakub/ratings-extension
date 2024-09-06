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

async function updateRating({ data, ratings }) {
  data.forEach(({ id, card }) => {
    if (card.querySelector(".mk-rating-container")) return;
    if (!ratings.has(id)) return;

    const { value } = ratings.get(id);
    const ratingContainer = createRatingContainer();
    ratingContainer.innerHTML = `${(+value).toFixed(2)}⭐️`;
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
  });
}

const data = new Map();

async function collectDataFromCards(cards) {
  for (let i = 0; i < cards.length; i++) {
    const card = cards[i];

    const title = card.getElementsByClassName("fallback-text")[0].innerText;
    const id = card
      .getElementsByTagName("a")[0]
      .href.split("/")[4]
      .split("?")[0];

    if (data.has(id)) continue;
    data.set(id, { title, id, card });
  }

  return { cards, data };
}

const ratings = new Map();
const queue = new Set();
let timeoutId = null;

async function fetchRatings({ cards, data }) {
  return new Promise((resolve) => {
    data.forEach(({ id }) => {
      if (ratings.has(id)) {
        return;
      }

      queue.add(id);

      // clear pending timeout
      if (timeoutId) clearTimeout(timeoutId);

      // create new timeout window
      timeoutId = setTimeout(() => {
        const query = Array.from(queue)
          .map((id) => {
            const res = data.get(id);
            return `${res.id}|${res.title}`;
          })
          .join("|");
        chrome.runtime.sendMessage(
          {
            type: "ratings",
            query,
          },
          (response) => {
            for (const r of response.ratings) {
              ratings.set(r.StreamingVendorId, { value: r.Value });
            }
          },
        );
        queue.clear();
        resolve({ cards, data, ratings });
      }, DELAY);
    });
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
