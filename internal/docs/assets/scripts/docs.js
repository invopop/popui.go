
document.addEventListener("alpine:init", () => {
  Alpine.data("docs", () => ({
    page: "getting-started",
    title: "Getting Started",
    init() {
      this.changePage();
      window.addEventListener("hashchange", () => {
        this.changePage();
      });
    },
    changePage() {
      const page = location.hash;
      if (page) {
        this.page = page.replace("#", "");
      }
      const main = document.querySelector("main");
      if (main) {
        main.scrollTo({ top: 0, behavior: "instant" });
      }
    },
    goto(id) {
      this.page = id;
    },
  }));

  Alpine.data("iframePreview", () => ({
    observer: null,
    init() {
      this.renderIframe();
      if (!this.observer) {
        this.observer = new MutationObserver(() => this.syncDark());
        this.observer.observe(document.documentElement, {
          attributes: true,
          attributeFilter: ["class"],
        });
      }
      this.$el.addEventListener("alpine:destroy", () => this.destroy(), {
        once: true,
      });
    },
    destroy() {
      if (this.observer) {
        this.observer.disconnect();
        this.observer = null;
      }
    },
    renderIframe() {
      const iframe = this.$refs.iframe;
      const doc = iframe.contentDocument || iframe.contentWindow.document;
      doc.open();
      doc.write(this.$refs.content.innerHTML);
      doc.close();
      this.syncDark();
    },
    syncDark() {
      const iframe = this.$refs.iframe;
      const doc = iframe.contentDocument || iframe.contentWindow.document;
      if (doc.documentElement) {
        doc.documentElement.classList.toggle(
          "dark",
          document.documentElement.classList.contains("dark"),
        );
      }
    },
  }));

  Alpine.data("darkModeToggle", () => ({
    dark: false,
    init() {
      const stored = localStorage.getItem("dark-mode");
      this.dark = stored !== null
        ? stored === "true"
        : document.documentElement.classList.contains("dark");
      document.documentElement.classList.toggle("dark", this.dark);
    },
    toggle() {
      this.dark = !this.dark;
      document.documentElement.classList.toggle("dark", this.dark);
      localStorage.setItem("dark-mode", this.dark);
    },
  }));

  Alpine.data("copyable", () => ({
    copied: false,
    copy(text) {
      navigator.clipboard.writeText(text);
      this.copied = true;
      setTimeout(() => (this.copied = false), 1500);
    },
  }));
});

