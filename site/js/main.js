/**
 * MoonBox Landing Page JavaScript
 * Adds interactivity and animations to the landing page
 */

document.addEventListener('DOMContentLoaded', function() {
    // Initialize variables
    const body = document.body;
    const themeToggleBtn = document.getElementById('theme-toggle-btn');
    const header = document.querySelector('.header');
    const navToggle = document.getElementById('nav-toggle');
    const navMenu = document.querySelector('.nav-menu');
    const navLinks = document.querySelectorAll('.nav-link');
    const backToTopBtn = document.getElementById('back-to-top');
    const copyButtons = document.querySelectorAll('.copy-btn');
    const langButtons = document.querySelectorAll('.lang-btn');

    // Check for saved theme preference or use user's system preference
    const savedTheme = localStorage.getItem('theme');
    const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');
    
    if (savedTheme === 'dark' || (!savedTheme && prefersDarkScheme.matches)) {
        body.setAttribute('data-theme', 'dark');
    }

    // Theme toggle functionality
    themeToggleBtn.addEventListener('click', function() {
        if (body.getAttribute('data-theme') === 'dark') {
            body.removeAttribute('data-theme');
            localStorage.setItem('theme', 'light');
        } else {
            body.setAttribute('data-theme', 'dark');
            localStorage.setItem('theme', 'dark');
        }
    });

    // Header scroll effect
    window.addEventListener('scroll', function() {
        if (window.scrollY > 50) {
            header.classList.add('scrolled');
            backToTopBtn.classList.add('visible');
        } else {
            header.classList.remove('scrolled');
            backToTopBtn.classList.remove('visible');
        }
    });

    // Mobile menu toggle
    navToggle.addEventListener('click', function() {
        navToggle.classList.toggle('active');
        navMenu.classList.toggle('active');
    });

    // Close mobile menu when clicking a link
    navLinks.forEach(link => {
        link.addEventListener('click', function() {
            navToggle.classList.remove('active');
            navMenu.classList.remove('active');
        });
    });

    // Smooth scrolling for anchor links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
            e.preventDefault();
            const targetId = this.getAttribute('href');
            if (targetId === '#') return;
            
            const targetElement = document.querySelector(targetId);
            if (targetElement) {
                window.scrollTo({
                    top: targetElement.offsetTop - 80,
                    behavior: 'smooth'
                });
            }
        });
    });

    // Back to top button
    backToTopBtn.addEventListener('click', function() {
        window.scrollTo({
            top: 0,
            behavior: 'smooth'
        });
    });

    // Copy code functionality
    copyButtons.forEach(button => {
        button.addEventListener('click', function() {
            const codeBlock = this.parentElement;
            const code = codeBlock.querySelector('code').innerText;
            
            navigator.clipboard.writeText(code).then(() => {
                // Show copied feedback
                const originalHTML = this.innerHTML;
                this.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>';
                
                setTimeout(() => {
                    this.innerHTML = originalHTML;
                }, 2000);
            });
        });
    });

    // Language switch
    langButtons.forEach(button => {
        button.addEventListener('click', function() {
            langButtons.forEach(btn => btn.classList.remove('active'));
            this.classList.add('active');
        });
    });

    // Initialize AOS (Animate on Scroll)
    initAOS();

    // Parallax effect for hero image
    const heroImage = document.querySelector('.main-image');
    const deploymentImage = document.querySelector('.parallax-image');
    
    if (heroImage) {
        initParallaxEffect(heroImage);
    }
    
    if (deploymentImage) {
        initParallaxEffect(deploymentImage);
    }

    // Active nav link based on scroll position
    updateActiveNavLink();
    window.addEventListener('scroll', updateActiveNavLink);
});

/**
 * Initialize Animate on Scroll effects
 * Simple implementation without external library
 */
function initAOS() {
    const animatedElements = document.querySelectorAll('[data-aos]');
    
    // Initial check for elements in viewport
    checkElementsInViewport(animatedElements);
    
    // Check on scroll
    window.addEventListener('scroll', function() {
        checkElementsInViewport(animatedElements);
    });
}

/**
 * Check if elements are in viewport and animate them
 * @param {NodeList} elements - Elements with data-aos attribute
 */
function checkElementsInViewport(elements) {
    elements.forEach(element => {
        const elementTop = element.getBoundingClientRect().top;
        const elementVisible = 150;
        
        if (elementTop < window.innerHeight - elementVisible) {
            element.classList.add('aos-animate');
        }
    });
}

/**
 * Initialize parallax effect for an image
 * @param {HTMLElement} element - The element to apply parallax effect to
 */
function initParallaxEffect(element) {
    document.addEventListener('mousemove', function(e) {
        const mouseX = e.clientX / window.innerWidth - 0.5;
        const mouseY = e.clientY / window.innerHeight - 0.5;
        
        element.style.transform = `perspective(1000px) rotateY(${mouseX * 10}deg) rotateX(${-mouseY * 10}deg)`;
    });
}

/**
 * Update active navigation link based on scroll position
 */
function updateActiveNavLink() {
    const sections = document.querySelectorAll('section[id]');
    const navLinks = document.querySelectorAll('.nav-link');
    
    let currentSection = '';
    
    sections.forEach(section => {
        const sectionTop = section.offsetTop - 100;
        const sectionHeight = section.offsetHeight;
        const sectionId = section.getAttribute('id');
        
        if (window.scrollY >= sectionTop && window.scrollY < sectionTop + sectionHeight) {
            currentSection = sectionId;
        }
    });
    
    navLinks.forEach(link => {
        link.classList.remove('active');
        if (link.getAttribute('href') === `#${currentSection}`) {
            link.classList.add('active');
        }
    });
}
