document.addEventListener('DOMContentLoaded', () => {
    // Initialize Feather icons
    feather.replace();

    // DOM Elements
    const dropZone = document.getElementById('dropZone');
    const fileInput = document.getElementById('fileInput');
    const browseBtn = document.getElementById('browseBtn');
    const fileStatus = document.getElementById('fileStatus');
    const fileName = document.getElementById('fileName');
    const removeFileBtn = document.getElementById('removeFileBtn');
    const patentContextSection = document.getElementById('patentContextSection');
    const tokenCount = document.getElementById('tokenCount');
    const contextPreview = document.getElementById('contextPreview');
    
    const chatForm = document.getElementById('chatForm');
    const messageInput = document.getElementById('messageInput');
    const chatMessages = document.getElementById('chatMessages');
    const sendBtn = document.getElementById('sendBtn');
    const clearChatBtn = document.getElementById('clearChatBtn');

    // State
    let patentContext = "";

    // -- File Upload Logic -- //

    // Open file dialog when browse is clicked
    browseBtn.addEventListener('click', () => fileInput.click());

    // Handle drag and drop events
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, preventDefaults, false);
    });

    function preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    ['dragenter', 'dragover'].forEach(eventName => {
        dropZone.addEventListener(eventName, () => dropZone.classList.add('drag-active'), false);
    });

    ['dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, () => dropZone.classList.remove('drag-active'), false);
    });

    dropZone.addEventListener('drop', (e) => {
        let dt = e.dataTransfer;
        let files = dt.files;
        if (files.length > 0) handleFile(files[0]);
    });

    fileInput.addEventListener('change', function() {
        if (this.files.length > 0) handleFile(this.files[0]);
    });

    function handleFile(file) {
        if (!file) return;

        fileName.textContent = file.name;
        dropZone.classList.add('hidden');
        fileStatus.classList.remove('hidden');

        const reader = new FileReader();
        reader.onload = (e) => {
            patentContext = e.target.result;
            
            // Rough token estimate (words / 0.75)
            const words = patentContext.trim().split(/\s+/).length;
            const estimatedTokens = Math.round(words / 0.75);
            
            tokenCount.textContent = `~${estimatedTokens.toLocaleString()} tokens`;
            contextPreview.textContent = patentContext.substring(0, 300) + (patentContext.length > 300 ? "..." : "");
            patentContextSection.classList.remove('hidden');
            
            addMessage("System", `Successfully loaded patent document: ${file.name}. You can now ask questions about it.`, "system");
        };
        reader.readAsText(file);
    }

    // Remove file
    removeFileBtn.addEventListener('click', () => {
        patentContext = "";
        fileInput.value = "";
        fileStatus.classList.add('hidden');
        patentContextSection.classList.add('hidden');
        dropZone.classList.remove('hidden');
        addMessage("System", "Patent document removed.", "system");
    });


    // -- Chat Logic -- //

    // Handle Enter to submit, Shift+Enter for newline
    messageInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            chatForm.dispatchEvent(new Event('submit'));
        }
    });

    // Auto-resize textarea
    messageInput.addEventListener('input', function() {
        this.style.height = 'auto';
        this.style.height = (this.scrollHeight) + 'px';
        if (this.value.trim() === '') {
            sendBtn.disabled = true;
        } else {
            sendBtn.disabled = false;
        }
    });

    chatForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const text = messageInput.value.trim();
        if (!text) return;

        // Reset input
        messageInput.value = '';
        messageInput.style.height = 'auto';
        sendBtn.disabled = true;

        // Add user message
        addMessage("You", text, "user");

        // Show loading indicator in a system message
        const loadingId = addMessage("LLM", "Thinking...", "system", true);

        try {
            // Call LLM endpoint
            const response = await fetch('http://localhost:8000/api/chat', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    message: text,
                    context: patentContext
                })
            });

            if (!response.ok) {
                throw new Error(`Server returned ${response.status} ${response.statusText}`);
            }

            const data = await response.json();
            
            // Assuming backend returns { response: "..." } or { text: "..." }
            const reply = data.response || data.text || data.message || JSON.stringify(data);
            
            updateMessage(loadingId, reply);
            
        } catch (error) {
            console.error("Chat API Error:", error);
            updateMessage(loadingId, `**Error:** Failed to connect to LLM. Make sure the backend is running at http://localhost:8000/api/chat.\n\nDetails: ${error.message}`);
        }
    });

    clearChatBtn.addEventListener('click', () => {
        chatMessages.innerHTML = '';
        addMessage("System", "Chat history cleared. How can I help you with your patent?", "system");
    });

    // Utility to add messages to the DOM
    function addMessage(sender, text, type = "system", isLoading = false) {
        const id = 'msg-' + Date.now();
        const msgDiv = document.createElement('div');
        msgDiv.className = `message ${type}`;
        msgDiv.id = id;

        const icon = type === 'user' ? 'user' : 'monitor';
        
        msgDiv.innerHTML = `
            <div class="avatar"><i data-feather="${icon}"></i></div>
            <div class="message-content ${isLoading ? 'loading' : ''}">${escapeHTML(text)}</div>
        `;
        
        chatMessages.appendChild(msgDiv);
        feather.replace(); // Render new icons
        
        // Scroll to bottom
        chatMessages.scrollTop = chatMessages.scrollHeight;
        
        return id;
    }

    // Utility to update an existing message (used for replacing "Thinking..." with actual text)
    function updateMessage(id, text) {
        const msgDiv = document.getElementById(id);
        if (msgDiv) {
            const contentDiv = msgDiv.querySelector('.message-content');
            contentDiv.classList.remove('loading');
            
            // Simple markdown parsing for bold and line breaks
            let formattedText = escapeHTML(text)
                .replace(/\n/g, '<br>')
                .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
                .replace(/`(.*?)`/g, '<code>$1</code>');
                
            contentDiv.innerHTML = formattedText;
            chatMessages.scrollTop = chatMessages.scrollHeight;
        }
    }

    // Basic HTML escaper
    function escapeHTML(str) {
        return str
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#039;');
    }
});
