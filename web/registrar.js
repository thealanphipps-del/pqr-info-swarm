document.addEventListener('DOMContentLoaded', () => {
    const searchBtn = document.getElementById('btn-search');
    const domainInput = document.getElementById('domain-search');
    const resultDiv = document.getElementById('search-result');
    
    let currentDomain = '';
    let currentPrice = 0;

    searchBtn.addEventListener('click', async () => {
        const domain = domainInput.value.trim();
        if (!domain) return;

        searchBtn.disabled = true;
        searchBtn.textContent = 'Searching...';

        try {
            const res = await fetch(`/REST/2.0/registrar/search?domain=${encodeURIComponent(domain)}`);
            const data = await res.json();

            resultDiv.classList.remove('d-none');
            
            if (data.available) {
                currentDomain = data.domain;
                currentPrice = data.price_usd;
                
                resultDiv.innerHTML = `
                    <div class="alert alert-success bg-transparent border-success text-success">
                        <h5>${data.domain} is available!</h5>
                        <p class="mb-2">Price: $${data.price_usd.toFixed(2)} ${data.currency} / year</p>
                        <button class="btn btn-outline-success" onclick="openCheckout()">Buy Now with Crypto</button>
                    </div>
                `;
            } else {
                resultDiv.innerHTML = `
                    <div class="alert alert-danger bg-transparent border-danger text-danger">
                        <h5>${data.domain} is taken.</h5>
                        <p>Try searching for a different domain.</p>
                    </div>
                `;
            }
        } catch (e) {
            resultDiv.innerHTML = `<div class="text-danger">Search failed: ${e.message}</div>`;
            resultDiv.classList.remove('d-none');
        } finally {
            searchBtn.disabled = false;
            searchBtn.textContent = 'Search';
        }
    });

    const purchaseBtn = document.getElementById('btn-purchase');
    purchaseBtn.addEventListener('click', async () => {
        const txHash = document.getElementById('tx-hash').value.trim();
        const cryptoMethod = document.getElementById('crypto-select').value;
        
        if (!txHash) {
            alert('Please enter a transaction hash.');
            return;
        }

        purchaseBtn.disabled = true;
        purchaseBtn.textContent = 'Verifying...';

        try {
            const res = await fetch('/REST/2.0/registrar/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    domain: currentDomain,
                    payment_method: cryptoMethod,
                    tx_hash: txHash
                })
            });

            const data = await res.json();
            
            if (res.ok) {
                alert(data.message);
                const modal = bootstrap.Modal.getInstance(document.getElementById('checkoutModal'));
                modal.hide();
                resultDiv.innerHTML = `<div class="alert alert-info">Domain ${currentDomain} has been registered!</div>`;
            } else {
                alert('Registration failed: ' + data.error);
            }
        } catch (e) {
            alert('Error: ' + e.message);
        } finally {
            purchaseBtn.disabled = false;
            purchaseBtn.textContent = 'Complete Purchase';
        }
    });
});

function openCheckout() {
    document.getElementById('checkout-domain').textContent = document.getElementById('search-result').querySelector('h5').innerText.split(' ')[0];
    document.getElementById('checkout-price').textContent = document.getElementById('search-result').querySelector('p').innerText.split(':')[1];
    
    const modal = new bootstrap.Modal(document.getElementById('checkoutModal'));
    modal.show();
}
