const options = [
  {
    value: "hideRatings",
    label: "Ukryj oceny",
    defaultValue: false,
  },
];
const optionsElem = document.getElementById("options");

options.forEach((option) => {
  const input = document.createElement("input");
  input.type = "checkbox";
  input.id = option.value;
  chrome.storage.local.get([option.value], (result) => {
    input.checked = result[option.value] ?? option.defaultValue;
  });
  input.addEventListener("change", () => {
    chrome.storage.local.set({ [option.value]: input.checked });
  });

  const label = document.createElement("label");
  label.htmlFor = option.value;
  label.textContent = option.label;

  optionsElem.appendChild(input);
  optionsElem.appendChild(label);
  optionsElem.appendChild(document.createElement("br"));
});
