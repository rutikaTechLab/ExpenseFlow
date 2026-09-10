let editingExpenseId = null;


// ========================================
// PAGE LOAD
// ========================================

document.addEventListener("DOMContentLoaded", function () {
    setToday();
    setupNavigation();
});


// ========================================
// SET TODAY
// ========================================

function setToday() {

    const dateInput = document.getElementById("expenseDate");

    if (dateInput && !dateInput.value) {

        const today = new Date();

        const year = today.getFullYear();
        const month = String(today.getMonth() + 1).padStart(2, "0");
        const day = String(today.getDate()).padStart(2, "0");

        dateInput.value = `${year}-${month}-${day}`;
    }
}


// ========================================
// EXPENSE MODAL
// ========================================

function openExpenseModal() {

    editingExpenseId = null;

    const modal = document.getElementById("expenseModal");
    const title = document.getElementById("expenseModalTitle");
    const form = document.getElementById("expenseForm");

    if (!modal || !title || !form) {
        console.error("Expense modal elements not found.");
        return;
    }

    title.textContent = "Add Expense";

    form.reset();

    setToday();

    modal.classList.add("show");

    setTimeout(function () {

        const description =
            document.getElementById("description");

        if (description) {
            description.focus();
        }

    }, 150);
}


// ========================================
// EDIT EXPENSE
// ========================================

function openEditModal(
    id,
    description,
    amount,
    category,
    date
) {

    editingExpenseId = id;

    const title =
        document.getElementById("expenseModalTitle");

    const descriptionInput =
        document.getElementById("description");

    const amountInput =
        document.getElementById("amount");

    const categoryInput =
        document.getElementById("category");

    const dateInput =
        document.getElementById("expenseDate");

    const modal =
        document.getElementById("expenseModal");

    if (
        !title ||
        !descriptionInput ||
        !amountInput ||
        !categoryInput ||
        !dateInput ||
        !modal
    ) {
        console.error("Edit modal elements not found.");
        return;
    }

    title.textContent = "Edit Expense";

    descriptionInput.value = description;
    amountInput.value = amount;
    categoryInput.value = category;
    dateInput.value = date;

    modal.classList.add("show");

    setTimeout(function () {
        descriptionInput.focus();
    }, 150);
}


// ========================================
// CLOSE EXPENSE MODAL
// ========================================

function closeExpenseModal() {

    const modal =
        document.getElementById("expenseModal");

    if (modal) {
        modal.classList.remove("show");
    }

    editingExpenseId = null;
}


// ========================================
// SAVE EXPENSE
// ========================================

async function saveExpense(event) {

    event.preventDefault();

    console.log("Save Expense clicked");


    // ------------------------------------
    // GET FORM VALUES
    // ------------------------------------

    const descriptionInput =
        document.getElementById("description");

    const amountInput =
        document.getElementById("amount");

    const categoryInput =
        document.getElementById("category");

    const dateInput =
        document.getElementById("expenseDate");


    if (
        !descriptionInput ||
        !amountInput ||
        !categoryInput ||
        !dateInput
    ) {

        console.error("One or more form fields are missing.");

        showToast(
            "Expense form is not configured correctly.",
            "error"
        );

        return;
    }


    const description =
        descriptionInput.value.trim();

    const amount =
        parseFloat(amountInput.value);

    const category =
        categoryInput.value;

    const date =
        dateInput.value;


    console.log("Expense data:", {
        description: description,
        amount: amount,
        category: category,
        date: date
    });


    // ------------------------------------
    // VALIDATION
    // ------------------------------------

    if (!description) {

        showToast(
            "Please enter a description.",
            "error"
        );

        descriptionInput.focus();

        return;
    }


    if (isNaN(amount) || amount <= 0) {

        showToast(
            "Please enter a valid amount.",
            "error"
        );

        amountInput.focus();

        return;
    }


    if (!category) {

        showToast(
            "Please select a category.",
            "error"
        );

        categoryInput.focus();

        return;
    }


    if (!date) {

        showToast(
            "Please select a date.",
            "error"
        );

        dateInput.focus();

        return;
    }


    // ------------------------------------
    // REQUEST DATA
    // ------------------------------------

    const expense = {
        description: description,
        amount: amount,
        category: category,
        date: date
    };


    let url = "/add-expense";
    let method = "POST";


    // ------------------------------------
    // EDIT MODE
    // ------------------------------------

    if (editingExpenseId !== null) {

        url =
            "/update-expense?id=" +
            encodeURIComponent(editingExpenseId);

        method = "PUT";
    }


    console.log("Sending request:", method, url);


    // ------------------------------------
    // SEND TO GO SERVER
    // ------------------------------------

    try {

        const response =
            await fetch(
                url,
                {
                    method: method,

                    headers: {
                        "Content-Type": "application/json"
                    },

                    body: JSON.stringify(expense)
                }
            );


        console.log(
            "Server response:",
            response.status,
            response.statusText
        );


        // --------------------------------
        // READ SERVER RESPONSE
        // --------------------------------

        const responseText =
            await response.text();

        console.log(
            "Server response body:",
            responseText
        );


        // --------------------------------
        // SERVER ERROR
        // --------------------------------

        if (!response.ok) {

            let errorMessage =
                responseText;

            try {

                const errorJSON =
                    JSON.parse(responseText);

                if (errorJSON.message) {
                    errorMessage =
                        errorJSON.message;
                }

            } catch (e) {
                // Response was plain text
            }


            throw new Error(
                `Server error (${response.status}): ${errorMessage}`
            );
        }


        // --------------------------------
        // SUCCESS RESPONSE
        // --------------------------------

        let result = {};

        try {

            result =
                JSON.parse(responseText);

        } catch (e) {

            console.warn(
                "Server returned non-JSON response."
            );
        }


        console.log(
            "Expense saved successfully:",
            result
        );


        showToast(
            result.message ||
            (
                editingExpenseId !== null
                    ? "Expense updated successfully."
                    : "Expense added successfully."
            ),
            "success"
        );


        closeExpenseModal();


        // --------------------------------
        // REFRESH DASHBOARD
        // --------------------------------

        setTimeout(function () {

            window.location.reload();

        }, 700);


    } catch (error) {

        console.error(
            "SAVE EXPENSE ERROR:",
            error
        );


        // IMPORTANT:
        // Show the REAL error instead of
        // hiding it behind "Something went wrong"

        showToast(
            error.message ||
            "Could not save expense.",
            "error"
        );
    }
}


// ========================================
// DELETE EXPENSE
// ========================================

async function deleteExpense(id) {

    const confirmed =
        confirm(
            "Are you sure you want to delete this expense?"
        );


    if (!confirmed) {
        return;
    }


    try {

        const response =
            await fetch(
                "/delete-expense?id=" +
                encodeURIComponent(id),

                {
                    method: "DELETE"
                }
            );


        const responseText =
            await response.text();


        if (!response.ok) {

            throw new Error(
                `Delete failed (${response.status}): ${responseText}`
            );
        }


        let result = {};

        try {
            result = JSON.parse(responseText);
        } catch (e) {
            console.warn("Delete response was not JSON.");
        }


        showToast(
            result.message ||
            "Expense deleted successfully.",
            "success"
        );


        setTimeout(function () {

            window.location.reload();

        }, 700);


    } catch (error) {

        console.error(
            "DELETE ERROR:",
            error
        );


        showToast(
            error.message ||
            "Could not delete expense.",
            "error"
        );
    }
}


// ========================================
// SEARCH + FILTER
// ========================================

function filterExpenses() {

    const searchInput =
        document.getElementById("searchInput");

    const categoryFilter =
        document.getElementById("categoryFilter");


    if (!searchInput || !categoryFilter) {
        return;
    }


    const search =
        searchInput.value
            .toLowerCase()
            .trim();


    const category =
        categoryFilter.value;


    const rows =
        document.querySelectorAll(
            ".expense-row"
        );


    rows.forEach(function (row) {

        const text =
            row.textContent
                .toLowerCase();


        const rowCategory =
            row.dataset.category ||
            "";


        const matchesSearch =
            text.includes(search);


        const matchesCategory =
            category === "all" ||
            rowCategory === category;


        if (
            matchesSearch &&
            matchesCategory
        ) {

            row.style.display = "";

        } else {

            row.style.display = "none";
        }

    });
}


// ========================================
// BUDGET MODAL
// ========================================

function openBudgetModal() {

    const modal =
        document.getElementById("budgetModal");

    const input =
        document.getElementById("budgetAmount");


    if (!modal) {
        return;
    }


    modal.classList.add("show");


    setTimeout(function () {

        if (input) {
            input.focus();
        }

    }, 150);
}


function closeBudgetModal() {

    const modal =
        document.getElementById("budgetModal");

    if (modal) {
        modal.classList.remove("show");
    }
}


// ========================================
// SAVE BUDGET
// ========================================

async function saveBudget(event) {

    event.preventDefault();


    const input =
        document.getElementById("budgetAmount");


    if (!input) {

        showToast(
            "Budget field not found.",
            "error"
        );

        return;
    }


    const budget =
        parseFloat(input.value);


    if (isNaN(budget) || budget <= 0) {

        showToast(
            "Enter a valid budget.",
            "error"
        );

        input.focus();

        return;
    }


    try {

        const response =
            await fetch(
                "/set-budget",
                {
                    method: "POST",

                    headers: {
                        "Content-Type":
                            "application/json"
                    },

                    body: JSON.stringify({
                        budget: budget
                    })
                }
            );


        const responseText =
            await response.text();


        if (!response.ok) {

            throw new Error(
                `Budget update failed (${response.status}): ${responseText}`
            );
        }


        let result = {};

        try {
            result =
                JSON.parse(responseText);
        } catch (e) {
            console.warn("Budget response was not JSON.");
        }


        showToast(
            result.message ||
            "Monthly budget updated.",
            "success"
        );


        closeBudgetModal();


        setTimeout(function () {

            window.location.reload();

        }, 700);


    } catch (error) {

        console.error(
            "BUDGET ERROR:",
            error
        );


        showToast(
            error.message ||
            "Could not update budget.",
            "error"
        );
    }
}


// ========================================
// TOAST
// ========================================

function showToast(
    message,
    type = "success"
) {

    const toast =
        document.getElementById("toast");

    const toastMessage =
        document.getElementById("toastMessage");


    if (!toast || !toastMessage) {

        console.error(
            "Toast elements not found:",
            message
        );

        alert(message);

        return;
    }


    toastMessage.textContent =
        message;


    toast.classList.remove(
        "success",
        "error",
        "show"
    );


    toast.classList.add(type);


    setTimeout(function () {

        toast.classList.add("show");

    }, 10);


    setTimeout(function () {

        toast.classList.remove("show");

    }, 4000);
}


// ========================================
// MODAL OUTSIDE CLICK
// ========================================

document.addEventListener(
    "click",
    function (event) {

        const expenseModal =
            document.getElementById(
                "expenseModal"
            );

        const budgetModal =
            document.getElementById(
                "budgetModal"
            );


        if (
            expenseModal &&
            event.target === expenseModal
        ) {

            closeExpenseModal();
        }


        if (
            budgetModal &&
            event.target === budgetModal
        ) {

            closeBudgetModal();
        }

    }
);


// ========================================
// ESCAPE KEY
// ========================================

document.addEventListener(
    "keydown",
    function (event) {

        if (event.key !== "Escape") {
            return;
        }


        closeExpenseModal();
        closeBudgetModal();

    }
);


// ========================================
// NAVIGATION
// ========================================

function setupNavigation() {

    const links =
        document.querySelectorAll(
            ".nav-link"
        );


    links.forEach(function (link) {

        link.addEventListener(
            "click",
            function () {

                links.forEach(
                    function (item) {

                        item.classList.remove(
                            "active"
                        );

                    }
                );


                link.classList.add(
                    "active"
                );

            }
        );

    });
}


// ========================================
// EDIT FROM TABLE ROW
// ========================================

function editExpenseFromRow(button) {

    const row =
        button.closest(".expense-row");


    if (!row) {

        console.error(
            "Expense row not found."
        );

        return;
    }


    const id =
        row.dataset.id;

    const description =
        row.dataset.description;

    const amount =
        row.dataset.amount;

    const category =
        row.dataset.category;

    const date =
        row.dataset.date;


    console.log(
        "Editing expense:",
        id,
        description,
        amount,
        category,
        date
    );


    openEditModal(
        id,
        description,
        amount,
        category,
        date
    );
}