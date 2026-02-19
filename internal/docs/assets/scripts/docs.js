// Debug HTMX events
document.addEventListener('htmx:beforeRequest', (e) => {
  console.log('HTMX beforeRequest:', e.detail);
});
document.addEventListener('htmx:afterRequest', (e) => {
  console.log('HTMX afterRequest:', e.detail);
});
document.addEventListener('htmx:beforeSwap', (e) => {
  console.log('HTMX beforeSwap:', e.detail);
});
document.addEventListener('htmx:afterSwap', (e) => {
  console.log('HTMX afterSwap:', e.detail);
});

document.addEventListener("alpine:init", () => {
  Alpine.data("docs", () => ({
    page: "getting-started",
    title: "Getting Started",
    init() {
      this.changePage();
      window.addEventListener("hashchange", () => {
        this.changePage();
      });
      // Process HTMX after Alpine initializes
      this.$nextTick(() => {
        if (typeof htmx !== 'undefined') {
          htmx.process(document.body);
        }
      });
    },
    changePage() {
      const page = location.hash;
      if (page) {
        this.page = page.replace("#", "");
      }
      // Scroll main area to top
      const main = document.querySelector('main');
      if (main) {
        main.scrollTo({ top: 0, behavior: 'instant' });
      }
      // Re-process HTMX when page changes
      this.$nextTick(() => {
        if (typeof htmx !== 'undefined') {
          htmx.process(document.body);
        }
      });
    },
    goto(id) {
      this.page = id;
    }
  }));
});
