/**
 * Подсветка активного пункта меню
 */
document.addEventListener('DOMContentLoaded', function () {
    const path = window.location.pathname;

    document.querySelectorAll('.nav-link').forEach(link => {
        const pattern = link.dataset.activePattern;

        if (pattern) {
            const regex = new RegExp(pattern);
            if (regex.test(path)) {
                link.classList.add('active');
            } else {
                link.classList.remove('active');
            }
        } else {
            // fallback для обычных ссылок
            const href = link.getAttribute('href');
            const isRoot = href === '/';
            const isMatch = isRoot ? path === '/' : path.startsWith(href);

            if (isMatch) {
                link.classList.add('active');
            } else {
                link.classList.remove('active');
            }
        }
    });
});
