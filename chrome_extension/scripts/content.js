const DELAY = 1_000;
const data = new Map();
const rendered = new Set();
let isDisabled = false;

// disableExtension
chrome.storage.local.get(["hideRatings"], (result) => {
  isDisabled = result.hideRatings;
});

chrome.storage.onChanged.addListener((changes) => {
  rendered.forEach(({ changeOpacity }) => {
    changeOpacity(changes.hideRatings.newValue ? 0 : 1);
  });
});

function createRatingContainer(href, value) {
  const ratingContainer = document.createElement("a");
  ratingContainer.href = href;
  ratingContainer.target = "_blank";
  ratingContainer.classList.add("mk-rating-container");
  const icon = document.createElement("img");
  const parsedValue = (+value).toFixed(2);
  let iconSrc = "assets/3.svg";
  switch (true) {
    case parsedValue >= 8:
      iconSrc = "assets/5.svg";
      break;
    case parsedValue >= 6:
      iconSrc = "assets/4.svg";
      break;
    case parsedValue >= 4:
      iconSrc = "assets/3.svg";
      break;
    case parsedValue >= 2:
      iconSrc = "assets/2.svg";
      break;
    case parsedValue > 0:
      iconSrc = "assets/1.svg";
      break;
  }

  icon.src = chrome.runtime.getURL(iconSrc);
  icon.style.width = "20px";
  icon.style.height = "20px";
  icon.style.marginRight = "5px";

  ratingContainer.appendChild(icon);
  ratingContainer.innerHTML += `${parsedValue}`;

  ratingContainer.style.top = "10px";
  ratingContainer.style.position = "absolute";
  ratingContainer.style.backgroundColor = "rgba(0, 0, 0, 0.5)";
  ratingContainer.style.color = "white";
  ratingContainer.style.padding = "3px";
  ratingContainer.style.borderRadius = "5px";
  ratingContainer.style.fontSize = "14px";
  ratingContainer.style.fontWeight = "bold";
  ratingContainer.style.verticalAlign = "middle";
  ratingContainer.style.display = "flex";
  ratingContainer.style.alignItems = "center";
  ratingContainer.style.justifyContent = "center";
  ratingContainer.style.transition = "opacity 0.3s";

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
    // const cards = document.getElementsByClassName("title-card");
    const cards = document.querySelectorAll(".title-card");
    // const modal = document.getElementsByClassName("previewModal--container");
    const modal = document.querySelectorAll(".previewModal--container");
    if (cards.length === 0) {
      await sleep(1000);
      return resolve(await getAllCards());
    }
    return resolve({ cards, modal });
  });
}

function getId(elem) {
  const aTags = elem.getElementsByTagName("a");
  if (aTags.length === 0) return;
  for (let i = 0; i < aTags.length; i++) {
    if (aTags[i].href.includes("title")) {
      // "title" seems to be the most reliable than "watch"
      return aTags[i].href.split("/")[4].split("?")[0];
    }
  }
  for (let i = 0; i < aTags.length; i++) {
    if (aTags[i].href.includes("watch")) {
      return aTags[i].href.split("/")[4].split("?")[0];
    }
  }
  return;
}

const ratingComponent = (value) => {
  const parsedValue = (+value).toFixed(2);
};

async function updateRating({ ratings, modal }) {
  data.forEach(({ id, cards }) => {
    cards.forEach((card) => {
      if (card.querySelector(".mk-rating-container")) return;
      if (!ratings.has(id)) return;

      const { value, href } = ratings.get(id);
      const ratingContainer = createRatingContainer(href, value);
      ratingContainer.style.right = "10px";
      ratingContainer.style.opacity = 0;
      setTimeout(() => {
        ratingContainer.style.opacity = isDisabled ? 0 : 1;
      }, DELAY);

      card.appendChild(ratingContainer);
      rendered.add({
        id,
        changeOpacity: (opacity) => {
          ratingContainer.style.opacity = opacity;
        },
      });
    });
  });

  if (modal[0] && !modal[0].querySelector(".mk-rating-container")) {
    const id = getId(modal[0]);
    if (!ratings.has(id)) return;
    const { value, href } = ratings.get(id);
    const ratingContainer = createRatingContainer(href, value);
    ratingContainer.style.left = "10px";
    modal[0].appendChild(ratingContainer);
  }
}

async function collectDataFromCards({ cards, modal }) {
  for (let i = 0; i < cards.length; i++) {
    const card = cards[i];

    const title =
      card?.getElementsByClassName("fallback-text")?.[0]?.innerText || "";
    const id = getId(card);

    if (data.has(id)) {
      data.get(id).cards.add(card);
    } else {
      data.set(id, { title, id, cards: new Set([card]) });
    }
  }

  return { cards, data, modal };
}

const ratings = new Map();
const queue = new Set();
let timeoutId = null;

async function fetchRatings({ cards, modal, data }) {
  return new Promise((resolve) => {
    data.forEach(({ id }) => {
      if (!ratings.has(id)) {
        queue.add(id);
      }

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
        queue.clear();

        chrome.runtime.sendMessage(
          {
            type: "ratings",
            query,
          },
          (response) => {
            for (const r of response.ratings) {
              ratings.set(r.StreamingVendorId, {
                value: r.Value,
                href: r.Link,
              });
            }
            resolve({ cards, ratings, modal });
          },
        );
      }, 500);
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
