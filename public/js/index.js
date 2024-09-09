/*
 * Collapse
 */

const MAX_PARAGRAPH_LENGTH = 420;

[...document.getElementsByClassName("expandable")].forEach((elm) => {
  const expandParagraph = document.createElement("p");
  const expandNewParagraph = document.createElement("p");

  const expandLink = document.createElement("a");
  expandLink.setAttribute("href", "javascript:void(0);")
  expandLink.textContent = "Read more";
  expandLink.onclick = () => {
    [...elm.children].forEach((child) => {
      child.style.display = 'block';
    })
    expandParagraph.style.display = 'none';
    expandNewParagraph.style.display = 'none';
  }

  [...elm.children].forEach((child) => {
    child.style.display = 'none'
  })

  if (elm.childElementCount > 1) {
    let text = elm.children[0].textContent;
    if (text.length > MAX_PARAGRAPH_LENGTH) {
      expandParagraph.append(text.substr(0, MAX_PARAGRAPH_LENGTH) + "... ");
      expandParagraph.appendChild(expandLink);
      expandParagraph.append(".");
      elm.appendChild(expandParagraph);
    } else {
      expandParagraph.append(text);
      expandNewParagraph.append("... ");
      expandNewParagraph.appendChild(expandLink);
      expandNewParagraph.append(".");
      elm.appendChild(expandParagraph);
      elm.appendChild(expandNewParagraph);
    }
  } else {
    let text = elm.children[0].textContent;
    if (text.length > MAX_PARAGRAPH_LENGTH) {
      expandParagraph.append(text.substr(0, MAX_PARAGRAPH_LENGTH) + "... ");
      expandParagraph.appendChild(expandLink);
      expandParagraph.append(".");
      elm.appendChild(expandParagraph);
    } else {
      expandParagraph.append(text);
      elm.appendChild(expandParagraph);
    }
  }
})
/*
 * Form
 */

const submitForm = (event) => {
  event.preventDefault();
  const form = event.target;
  const formData = new FormData(form);
  const xhr = new XMLHttpRequest();
  xhr.addEventListener("load", () => {
      if (xhr.status >= 400) {
        let targetErrorQuery = event.target.getAttribute("data-target-error");
        let targetError = document.querySelector(targetErrorQuery);
        if (targetError) {
          targetError.setAttribute("aria-invalid", "true");
          decribedById = targetError.getAttribute("aria-describedby");
          let describedBy = document.getElementById(decribedById);
          describedBy.textContent = xhr.responseText;
        }
      } else {
        location.reload();
      }
  });
  xhr.open("POST", form.action, true);
  xhr.send(formData);  
}

/*
 * Modal
 */

// Config
const isOpenClass = "modal-is-open";
const scrollbarWidthCssVar = "--pico-scrollbar-width";
let visibleModal = null;

// Toggle modal
const toggleModal = (event) => {
  event.preventDefault();
  const modal = document.getElementById(event.currentTarget.dataset.target);
  if (!modal) return;
  modal && (modal.open ? closeModal(modal) : openModal(modal));
};

// Open modal
const openModal = (modal) => {
  const { documentElement: html } = document;
  const scrollbarWidth = getScrollbarWidth();
  if (scrollbarWidth) {
    html.style.setProperty(scrollbarWidthCssVar, `${scrollbarWidth}px`);
  }
  html.classList.add(isOpenClass);
  // timeout to not trigger the outside-of-modal click event
  // that closes it.
  setTimeout(() => { visibleModal = modal; }, 100 /* ms */);
  modal.showModal();
};

// Close modal
const closeModal = (modal) => {
  visibleModal = null;
  const { documentElement: html } = document;
  html.classList.remove(isOpenClass);
  html.style.removeProperty(scrollbarWidthCssVar);
  modal.close();
};

// Close with a click outside
document.addEventListener("click", (event) => {
  if (visibleModal === null) return;
  const modalContent = visibleModal.querySelector("article");
  const isClickInside = modalContent.contains(event.target);
  !isClickInside && closeModal(visibleModal);
});

// Close with Esc key
document.addEventListener("keydown", (event) => {
  if (event.key === "Escape" && visibleModal) {
    closeModal(visibleModal);
  }
});

// Get scrollbar width
const getScrollbarWidth = () => {
  const scrollbarWidth = window.innerWidth - document.documentElement.clientWidth;
  return scrollbarWidth;
};

// Is scrollbar visible
const isScrollbarVisible = () => {
  return document.body.scrollHeight > screen.height;
};
