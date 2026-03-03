
document.addEventListener("alpine:init", () => {
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

  Alpine.data("copyable", (text) => ({
    copied: false,
    copy() {
      navigator.clipboard.writeText(text);
      this.copied = true;
      setTimeout(() => (this.copied = false), 1500);
    },
  }));
});
