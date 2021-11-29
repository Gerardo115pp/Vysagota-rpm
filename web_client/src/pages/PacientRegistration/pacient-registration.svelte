<script>
    import Header from "../../components/Header.svelte";
    import Input from "../../components/Input.svelte";
    import Button from "../../components/Button.svelte";
    import FieldData, { FieldStates, verifyFormFields } from "../../classes/FieldData";
    import Pacient from "../../classes/Pacient";
    import { push } from "svelte-spa-router";


    let is_form_ready = false;
    const pacient = new Pacient();

    let form_data = [
        new FieldData("pacient_name", /^([a-zA-Z]+\s?){1,3}$/, "pacient"),
        new FieldData("pacient_username", /^[a-zA-Z_\-\d]+$/, "username"),
        new FieldData("pacient_password", /^[A-z_\-+\d\$%\^@!&#\.<>\?]+$/, "password", "password"),
        new FieldData("pacient_email", /^[a-zA-Z_\d\.]+@(\.?[a-z\d]+){2,10}$/, "email"),
        new FieldData("pacient_email_tutor",/^[a-zA-Z_\d\.]+@(\.?[a-z\d]+){2,10}$/, "tutor_email"),
        new FieldData("pacient_phone", /^\d{9,12}$/, "phone"),
        new FieldData("pacient_birth", /^\d{4}-\d{2}-\d{2}$/, "birth"),
        new FieldData("pacient_gender", /^([Mm]ale|[Ff]emale)$/, "gender", "text",false)
    ]

    form_data[6].setPlaceholder("Birthday (YYYY-MM-DD)");

    // form data link with Pacient obj
        $: pacient.name = form_data[0].getFieldValue();
        $: pacient.username = form_data[1].getFieldValue();
        $: pacient.secret = form_data[2].getFieldValue();
        $: pacient.email = form_data[3].getFieldValue();
        $: pacient.tutor_email = form_data[4].getFieldValue();
        $: pacient.phone = form_data[5].getFieldValue();
        $: pacient.birthday = form_data[6].isReady() ? new Date(form_data[6].getFieldValue()) : null;
        $: pacient.is_male = form_data[7].getFieldValue().toLowerCase() === "male" ? true : false;
    // end of data linking



    const verifyRegistrationForm = () => {
        let is_valid = verifyFormFields(form_data);
        is_form_ready = is_valid;
        form_data = [...form_data]; // trigger svelte update
    }

    const createPacient = () => {
        if (is_form_ready) {
            pacient.create();
            push("/doctor")
        }
    }
</script>

<div id="pacient-registration" class="page">
    <Header />
    <div id="pr-form">
        <div id="prf-header-label">
            Pacient Registration
        </div>
        <div id="form-fields">
            {#each form_data as field}
                <Input onBlur={verifyRegistrationForm} onEnterPressed={verifyRegistrationForm} field_data={field} isClear={true} isSquared={true}/>
            {/each}
        </div>
        <div id="controls">
            <Button isEnabled={is_form_ready} onClick={createPacient} label="Register" width="60%" padding="1.2vh 2vw" font_size="1.3em" isSquared={true}/>
        </div>
    </div>
</div>

<style>
    #pr-form {
        display: flex;
        height: 89vh;
        flex-direction: column;
        justify-content: space-around;
        align-items: center;
        /* padding: 15vh 0; */
    }

    #prf-header-label {
        font-family: var(--title-font);
        font-size: 2.5em;
    }

    #form-fields {
        display: flex;
        flex-direction: column;
        justify-content: space-around;
        align-items: center;
        width: 40%;
        height: 60vh;
    }
</style>