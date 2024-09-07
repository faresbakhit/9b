/*
 * Form
 */

const submitForm = (event) => {
  event.preventDefault();
  const form = event.target;
  const formData = new FormData(form);
  const xhr = new XMLHttpRequest();
  xhr.open("POST", form.action, true);
  xhr.onload = () => {
      if (xhr.status >= 400) {
        let targetError = document.querySelector(event.target.getAttribute("target-error"));
        targetError.setAttribute("aria-invalid", "true");
        let describedBy = document.getElementById(targetError.getAttribute("aria-describedby"));
        describedBy.textContent = xhr.responseText;
      } else {
        location.reload();
      }
  };
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
