<script>
    import Header from '../../components/Header.svelte';
    import Button from '../../components/Button.svelte';
    import { push } from 'svelte-spa-router';
    import { onMount } from "svelte";
    import { getPacients } from "../../classes/Pacient";

    let pacients = [];
    let active_pacient_index = 0;

    onMount(() => {
        getPacients(data => {
            console.log(data);
            pacients = data;
        });
    });

    const changeActivePacient = (index) => {
        active_pacient_index = index;
        console.log(pacients[index]);
    };

</script>

<div id="pacients-summary" class="page">
    <Header />
    <main id="page-content">
        <div id="main-content">
            <div id="pacients-controls">
                <div id="pc-buttons">
                    <div id="create-pacient-btn">
                        <Button onClick={() => push("/pacient-registration")}  width="100%" padding="1vh .5vw" label="+  Add Pacient" isSquared={true} />
                    </div>
                </div>
                <div class="hs"></div>
            </div>
            <div id="pacient-data">
                <div id="pd-header">
                    {active_pacient_index < pacients.length ? pacients[active_pacient_index].username : ""}
                </div>
                <div id="pd-container">
                    <div class="pdc-field">
                        <div class="pdcf-label">name</div>
                        <div class="pdcf-value">{active_pacient_index < pacients.length ? pacients[active_pacient_index].name : ""}</div>
                    </div>
                    <div class="pdc-field">
                        <div class="pdcf-label">Gender</div>
                        <div class="pdcf-value">{active_pacient_index < pacients.length ? pacients[active_pacient_index].gender : ""}</div>
                    </div>
                    <div class="pdc-field">
                        <div class="pdcf-label">Tutor email</div>
                        <div class="pdcf-value">{active_pacient_index < pacients.length ? pacients[active_pacient_index].tutor_email : ""}</div>
                    </div>
                    <div class="pdc-field">
                        <div class="pdcf-label">email</div>
                        <div class="pdcf-value">{active_pacient_index < pacients.length ? pacients[active_pacient_index].email : ""}</div>
                    </div>
                    <div class="pdc-field">
                        <div class="pdcf-label">appointments</div>
                        <div class="pdcf-value">{active_pacient_index < pacients.length ? pacients[active_pacient_index].appointments_count: ""}</div>
                    </div>
                </div>
            </div>
        </div>
        <div id="pacients-list">
            {#each pacients as pacient, index}
                <div on:click={() => changeActivePacient(index)} class="pacient-item">
                    <div class="pi-name">{pacient.name}</div>
                </div>
            {/each}
        </div>
    </main>
</div> 

<style>

    #page-content {
        display: flex;
        height: 90.5vh;
    }
    #main-content {
        width: 80%;
    }

    
    #pc-buttons {
        display: flex;
        width: 100%;
        height: 10vh;
        justify-content: flex-start;
        align-items: center;
        padding: 0 2vw;
    }

    #create-pacient-btn {
        width: 8%;
    }
  
    #pacient-data {
        height: 70.5vh;
    }
    
    #pd-header {
        height: 5vh;
        color: var(--dark-color);
        font-size: 1.2em;
        font-weight: 500;
        padding: 2vh 2vw;
    }

    #pd-container {
        box-sizing: border-box;
        display: flex;
        flex-direction: column;
        justify-content: space-around;
        height: 70vh;
        padding: 2vh 4vw;
    }

    .pdc-field {
        display: flex;
        justify-content: flex-start;
        height: 3vh;
    }

    .pdcf-label {
        font-weight: bold;
        margin-right: 1vw;
    }
    
    .pdcf-label::after {
        content: ":";
    }

    #pacients-list {
        font-family: var(--title-font);
        overflow-x: hidden;
        width: 20%;
        height: 88.5vh;
        padding: 0.5em;
        border-left: 1px solid var(--theme-color);
    }

    .pacient-item {
        cursor: pointer;
        font-size: 1.2em;
        padding: 2rem 2em;
        border-bottom: 1px solid var(--theme-color);
        transition: all .2s ease-in;
    }

    .pacient-item:last-child {
        border-bottom: none;
    }

    .pacient-item:hover {
        background-color: var(--theme-color);
        color: white;
    }
    
</style>