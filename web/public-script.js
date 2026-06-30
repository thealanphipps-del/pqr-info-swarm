// public-script.js

// --- 1. Custom Cursor ---
const cursorDot = document.querySelector('.cursor-dot');
const cursorOutline = document.querySelector('.cursor-outline');

let cursorScale = 1;

let cursorX = window.innerWidth / 2;
let cursorY = window.innerHeight / 2;
let outlineX = cursorX;
let outlineY = cursorY;

let mouseX = 0;
let mouseY = 0;

window.addEventListener('mousemove', (e) => {
    cursorX = e.clientX;
    cursorY = e.clientY;
    mouseX = e.clientX / window.innerWidth - 0.5;
    mouseY = e.clientY / window.innerHeight - 0.5;
});

function animateCursor() {
    cursorDot.style.transform = `translate3d(${cursorX}px, ${cursorY}px, 0) translate(-50%, -50%)`;
    outlineX += (cursorX - outlineX) * 0.15;
    outlineY += (cursorY - outlineY) * 0.15;
    cursorOutline.style.transform = `translate3d(${outlineX}px, ${outlineY}px, 0) translate(-50%, -50%) scale(${cursorScale})`;
    requestAnimationFrame(animateCursor);
}
animateCursor();

// --- 2. Magnetic Elements ---
const magnetics = document.querySelectorAll('.magnetic');

magnetics.forEach((magnet) => {
    magnet.addEventListener('mousemove', (e) => {
        const position = magnet.getBoundingClientRect();
        const x = e.clientX - position.left - position.width / 2;
        const y = e.clientY - position.top - position.height / 2;
        
        const strength = magnet.getAttribute('data-strength') || 20;

        gsap.to(magnet, {
            x: x * (strength / 100),
            y: y * (strength / 100),
            duration: 0.5,
            ease: "power2.out"
        });

        // Expand cursor outline on hover
        cursorScale = 1.5;
        cursorOutline.style.backgroundColor = 'rgba(0, 229, 255, 0.1)';
    });

    magnet.addEventListener('mouseleave', (e) => {
        gsap.to(magnet, {
            x: 0,
            y: 0,
            duration: 0.5,
            ease: "power2.out"
        });
        cursorScale = 1;
        cursorOutline.style.backgroundColor = 'transparent';
    });
});

// --- 3. Gate 6 Friction Overlay ---
const frictionBtn = document.getElementById('friction-btn');
const fillBar = document.querySelector('.fill-bar');
const overlay = document.getElementById('friction-overlay');
const mainWrapper = document.getElementById('main-wrapper');

let holdTimer;
let progress = 0;

frictionBtn.addEventListener('mousedown', startHold);
frictionBtn.addEventListener('touchstart', startHold, { passive: false });

frictionBtn.addEventListener('mouseup', endHold);
frictionBtn.addEventListener('mouseleave', endHold);
frictionBtn.addEventListener('touchend', endHold);
frictionBtn.addEventListener('touchcancel', endHold);

function startHold(e) {
    if (e && e.type === 'touchstart') e.preventDefault();
    if (holdTimer) clearInterval(holdTimer);
    holdTimer = setInterval(() => {
        progress += 2;
        fillBar.style.width = `${progress}%`;
        
        if (progress >= 100) {
            clearInterval(holdTimer);
            unlockMainSite();
        }
    }, 20);
}

function endHold() {
    clearInterval(holdTimer);
    if (progress < 100) {
        progress = 0;
        fillBar.style.width = `0%`;
    }
}

function unlockMainSite() {
    gsap.to(overlay, {
        opacity: 0,
        duration: 1,
        onComplete: () => {
            overlay.style.display = 'none';
            mainWrapper.style.visibility = 'visible';
            mainWrapper.style.opacity = 1;
            initScrollAnimations();
        }
    });
}

// --- 4. GSAP ScrollTrigger Animations ---
gsap.registerPlugin(ScrollTrigger);

function initScrollAnimations() {
    // Hero animations
    gsap.to('.fade-up', {
        y: 0,
        opacity: 1,
        duration: 1,
        stagger: 0.2,
        ease: "power3.out"
    });

    // Scroll elements
    const fadeUps = document.querySelectorAll('.section .fade-up:not(.hero .fade-up)');
    fadeUps.forEach((el) => {
        gsap.to(el, {
            scrollTrigger: {
                trigger: el,
                start: "top 85%",
            },
            y: 0,
            opacity: 1,
            duration: 1,
            ease: "power3.out"
        });
    });

    gsap.to('.fade-right', {
        scrollTrigger: { trigger: '.split-layout', start: "top 80%" },
        x: 0, opacity: 1, duration: 1, ease: "power3.out"
    });

    gsap.to('.fade-left', {
        scrollTrigger: { trigger: '.split-layout', start: "top 80%" },
        x: 0, opacity: 1, duration: 1, ease: "power3.out"
    });
}

// --- 5. Three.js WebGL Background (Triple Helix) ---
const canvas = document.getElementById('webgl-canvas');
const scene = new THREE.Scene();
const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
const renderer = new THREE.WebGLRenderer({ canvas: canvas, alpha: true, antialias: true });

renderer.setSize(window.innerWidth, window.innerHeight);
renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

// Create Particles
const particlesGeometry = new THREE.BufferGeometry();
const particlesCount = 2000;
const posArray = new Float32Array(particlesCount * 3);

for(let i = 0; i < particlesCount; i++) {
    const strand = i % 3;
    const t = (i / particlesCount) * Math.PI * 20;
    const offset = strand * ((Math.PI * 2) / 3);
    const radius = 2;
    
    const x = Math.cos(t + offset) * radius + (Math.random() - 0.5) * 0.3;
    const y = (i / particlesCount) * 10 - 5;
    const z = Math.sin(t + offset) * radius + (Math.random() - 0.5) * 0.3;

    posArray[i * 3] = x;
    posArray[i * 3 + 1] = y;
    posArray[i * 3 + 2] = z;
}

particlesGeometry.setAttribute('position', new THREE.BufferAttribute(posArray, 3));
const particlesMaterial = new THREE.PointsMaterial({
    size: 0.015,
    color: 0x00e5ff,
    transparent: true,
    opacity: 0.8,
    blending: THREE.AdditiveBlending
});

const particlesMesh = new THREE.Points(particlesGeometry, particlesMaterial);
scene.add(particlesMesh);
camera.position.z = 3;

// Mouse Interaction
// Mouse Interation State hoisted to top

window.addEventListener('resize', () => {
    camera.aspect = window.innerWidth / window.innerHeight;
    camera.updateProjectionMatrix();
    renderer.setSize(window.innerWidth, window.innerHeight);
});

const clock = new THREE.Clock();

function animate() {
    requestAnimationFrame(animate);
    const elapsedTime = clock.getElapsedTime();

    particlesMesh.rotation.y = elapsedTime * 0.1;
    particlesMesh.rotation.x = elapsedTime * 0.05;

    // Subtle parallax
    camera.position.x += (mouseX * 0.5 - camera.position.x) * 0.05;
    camera.position.y += (-mouseY * 0.5 - camera.position.y) * 0.05;

    renderer.render(scene, camera);
}

animate();
