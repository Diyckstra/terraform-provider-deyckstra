document.addEventListener("DOMContentLoaded", function () {
    console.log("Script is running...");

    // Перекраска всех <code> элементов по всему документу
    document.querySelectorAll('code').forEach(function (codeElement) {
        console.log("Found <code> element:", codeElement.innerHTML);

        if (codeElement.closest('a')) {
            console.log("Setting to blue:", codeElement.innerHTML);
            codeElement.classList.add('text-blue');
        } else {
            console.log("Setting to red:", codeElement.innerHTML);
            codeElement.classList.add('text-red');
        }
    });

    document.querySelectorAll('p').forEach(function (paragraph) {
        console.log("Processing paragraph:", paragraph);

        let iconHTML = '';
        let htmlContent = paragraph.innerHTML;

        // Обработка начального блока на note-block стиль
        if (htmlContent.startsWith('~&gt;')) {
            htmlContent = htmlContent.replace('~&gt;', '').trim();
            paragraph.classList.add('note-block', 'note-block-yellow');
            iconHTML = '<i class="fas fa-exclamation-triangle" style="color: #ffc107; margin-right: 0.5em;"></i>';
            console.log("Yellow note block detected.");

        } else if (htmlContent.startsWith('-&gt;')) {
            htmlContent = htmlContent.replace('-&gt;', '').trim();
            paragraph.classList.add('note-block', 'note-block-blue');
            iconHTML = '<i class="fas fa-info-circle" style="color: #007acc; margin-right: 0.5em;"></i>';
            console.log("Blue note block detected.");

        } else if (htmlContent.startsWith('!&gt;')) {
            htmlContent = htmlContent.replace('!&gt;', '').trim();
            paragraph.classList.add('note-block', 'note-block-red');
            iconHTML = '<i class="fas fa-exclamation-circle" style="color: #d9534f; margin-right: 0.5em;"></i>';
            console.log("Red note block detected.");
        }

        // Извлекаем <strong> элементы на 2 строки вверх и удаляем из текущего местоположения
        let strongText = '';
        const strongElements = paragraph.querySelectorAll('strong');
        
        strongElements.forEach(function (strongElement) {
            console.log("Found <strong> element:", strongElement.innerHTML);

            if (paragraph.classList.contains('note-block-yellow')) {
                strongElement.classList.add('text-yellow');
            } else if (paragraph.classList.contains('note-block-blue')) {
                strongElement.classList.add('text-blue');
            } else if (paragraph.classList.contains('note-block-red')) {
                strongElement.classList.add('text-red');
            }
            strongText += strongElement.outerHTML;
            strongElement.remove(); // Удаляем strongElement после обработки
        });

        // Обработанный контент создается только если есть <strong> элементы
        let finalContent = iconHTML + strongText;
        if (strongText) {
            finalContent += '<br><br>';
        }
        finalContent += htmlContent.trim().replace(/<strong>.*?<\/strong>/g, '');
        
        // Обновляем paragraph.innerHTML для вывода
        paragraph.innerHTML = finalContent;
        console.log("Updated HTML content:", paragraph.innerHTML);
    });
});
