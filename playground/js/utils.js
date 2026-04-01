/* ── DOM helpers ── */
const $ = (id) => document.getElementById(id);
const $$ = (sel) => document.querySelectorAll(sel);

/* ── Sleep ── */
function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

/* ── Timestamp ── */
function timestamp() {
  const d = new Date();
  return [d.getHours(), d.getMinutes(), d.getSeconds()]
    .map((n) => String(n).padStart(2, '0'))
    .join(':');
}

/* ── Copy to clipboard ── */
function copyCode(btn) {
  const container = btn.closest('.hero-code, .tab-panel, .code-block');
  const pre = container ? container.querySelector('pre') : btn.parentElement.nextElementSibling;
  if (!pre) return;
  navigator.clipboard.writeText(pre.textContent).then(() => {
    btn.textContent = 'Copied!';
    btn.classList.add('copied');
    setTimeout(() => {
      btn.textContent = 'Copy';
      btn.classList.remove('copied');
    }, 2000);
  });
}

/* ── Scroll reveal observer ── */
function initScrollReveal() {
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible');
          observer.unobserve(entry.target);
        }
      });
    },
    { threshold: 0.1, rootMargin: '0px 0px -40px 0px' }
  );
  document.querySelectorAll('.sr').forEach((el) => observer.observe(el));
}

/* ── Animated counter ── */
function animateValue(el, start, end, duration) {
  const range = end - start;
  const startTime = performance.now();
  function tick(now) {
    const elapsed = now - startTime;
    const progress = Math.min(elapsed / duration, 1);
    const eased = 1 - Math.pow(1 - progress, 3);
    el.textContent = Math.round(start + range * eased).toLocaleString();
    if (progress < 1) requestAnimationFrame(tick);
  }
  requestAnimationFrame(tick);
}
