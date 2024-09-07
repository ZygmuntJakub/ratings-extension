const options = [
  {
    value: "disableExtension",
    label: "Wyłącz rozszerzenie",
    defaultValue: false,
  },
];
const optionsElem = document.getElementById("options");

options.forEach((option) => {
  const input = document.createElement("input");
  input.type = "checkbox";
  input.id = option.value;
  input.checked = option.defaultValue;
  input.addEventListener("change", () => {
    chrome.storage.sync.set({ [option.value]: input.checked });
  });

  const label = document.createElement("label");
  label.htmlFor = option.value;
  label.textContent = option.label;

  optionsElem.appendChild(input);
  optionsElem.appendChild(label);
  optionsElem.appendChild(document.createElement("br"));
});
