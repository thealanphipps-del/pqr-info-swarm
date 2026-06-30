const canvas = document.getElementById('hudCanvas');
const ctx = canvas.getContext('2d');
let time = 0;

function resize() {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
}

window.addEventListener('resize', resize);
resize();

function drawHUD() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    // Draw some cool futuristic HUD elements
    ctx.strokeStyle = 'rgba(0, 229, 255, 0.5)';
    ctx.lineWidth = 2;
    
    ctx.save();
    ctx.translate(canvas.width / 2, canvas.height / 2);
    ctx.rotate(time * 0.005);
    
    // Outer circle
    ctx.beginPath();
    ctx.arc(0, 0, 200, 0, Math.PI * 2);
    ctx.stroke();
    
    // Inner dashed circle
    ctx.beginPath();
    ctx.setLineDash([10, 15]);
    ctx.arc(0, 0, 150, 0, Math.PI * 2);
    ctx.stroke();
    
    // Triple helix points
    ctx.fillStyle = '#00e5ff';
    for(let i=0; i<3; i++) {
        let angle = (i * Math.PI * 2 / 3) + (time * 0.02);
        let x = Math.cos(angle) * 100;
        let y = Math.sin(angle) * 100;
        ctx.beginPath();
        ctx.arc(x, y, 5, 0, Math.PI*2);
        ctx.fill();
    }
    
    ctx.restore();
    
    // Text elements
    ctx.fillStyle = '#00e5ff';
    ctx.font = '16px monospace';
    ctx.fillText('SYS.STATUS: ONLINE', 20, 30);
    ctx.fillText('LOC: ' + Math.floor(Math.random()*1000) + ' / ' + Math.floor(Math.random()*1000), 20, 50);
    
    time++;
    requestAnimationFrame(drawHUD);
}

drawHUD();

// Secret Trigger Logic
// Double clicking on the canvas will unlock the public site
canvas.addEventListener('dblclick', () => {
    window.location.href = 'public.html';
});
