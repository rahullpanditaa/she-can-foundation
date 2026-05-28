const form = document.getElementById("contact-form");
const successMessage = document.getElementById("success-message");
const submitButton = document.getElementById("submit-btn");
const toast = document.getElementById("toast");

form.addEventListener("submit", async (event) => {
    event.preventDefault();
    submitButton.disabled = true;
    submitButton.innerText = "Submitting...";

    const formData = {
        name: document.getElementById("name").value,
        email: document.getElementById("email").value,
        message: document.getElementById("message").value,
    };

    // if name or email or message not present
    if (!formData.name || !formData.email || !formData.message) {
        successMessage.innerText = "Please fill all fields";
        submitButton.disabled = false;
        submitButton.innerText = "Submit";
        return;
    }

    try {
        const response = await fetch("https://she-can-foundation-czjf.onrender.com/submit", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(formData),
        });

        const data = await response.json();
        if (data.success) {
            toast.innerText = "Form Submitted Successfully";
            toast.classList.add("show");

            setTimeout(() => {
                toast.classList.remove("show");
            }, 3000);

            form.reset();
        }
    } catch (error) {
        successMessage.innerText = "Something went wrong";
    }

    submitButton.disabled = false;
    submitButton.innerText = "Submit";
});