<script>
    import Header from '../../components/Header.svelte';
    import Input from '../../components/Input.svelte';
    import Button from '../../components/Button.svelte';
    import FieldData, { verifyFormFields } from '../../classes/FieldData';
    import { push } from 'svelte-spa-router';
    import vysagota_client from '../../classes/VysagotaClient';
    
    export let params = {};
    
    // DELETE: no longer needed
    const profile_type = params['profile'] != undefined ? params['profile'] : 'user';
    
    const form_data = [
        new FieldData("pacient_username", /^[a-zA-Z\-_\d]+$/, "username", "text", true),
        new FieldData("pacient_password", /^[a-zA-Z_\d]+$/, "password", "password", true)
    ]
    
    let is_form_ready = false;

    const redirect = (retry=true) => {
        if (is_form_ready) {
            let service_data = vysagota_client.getEndpointData('authorization');
            if (service_data == undefined) {
                vysagota_client.getEndpoints(() => {
                    redirect(false);
                });
                return;
            }
            console.log(service_data);
            const user_data = new FormData();
            user_data.append('username', form_data[0].getFieldValue());
            user_data.append('password', form_data[1].getFieldValue());
            
            const request = new Request(`http://${service_data.host}:${service_data.port}/login`, { method: 'POST', body: user_data});
            fetch(request)
            .then(promise => promise.json())
            .then(response => {
                if (response.status) { 
                    push(`/${response.type}`);
                }
            })
        }
    }
    
    const isFormReady = () => {
        let ready = verifyFormFields(form_data);
        is_form_ready = ready;
    }

</script>

<div id="login" class="page">
    <Header />
    <div class="subtitle">
        Welcome back, <id id="user-name">{profile_type}</id>
    </div>
    <div id="login-form">
        <div id="login-header-label">
            Login 
        </div>
        <div id="form-fields">
            {#each form_data as field}
                <Input onBlur={isFormReady} onEnterPressed={redirect} field_data={field} isSquared={true}/>
            {/each}
        </div>
        <div id="controls">
            <Button onClick={redirect} isEnabled={is_form_ready} label="Login" width="30%" padding="1.2vh 2vw" font_size="1.3em" isSquared={true}/>
        </div>
    </div>
</div>

<style>
    .subtitle {
        font-size: 0.9em;
        font-family: var(--title-font);
        padding: .5em 0 0 2em;
    }

    #login-form {
        display: flex;
        height: 80vh;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }

    #login-header-label {
        display: flex;
        height: 10vh;
        font-family: var(--title-font);
        font-size: 2.5em;
    }

    #form-fields {
        display: flex;
        width: 20%;
        height: 15vh;
        flex-direction: column;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 2vh;
    }

    #controls {
        display: flex;
        width: 20%;
        height: 10vh;
        flex-direction: column;
        align-items: center;
        justify-content: space-between;
    }
</style>