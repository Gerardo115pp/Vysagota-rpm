<script>
    import Header from "../../components/Header.svelte";
    import { push } from "svelte-spa-router";
</script>

<div id="doctor-hub">
    <Header return_to="/" />
    <div id="hub-content">
        <div id="doctor-hub-entrys">
            <div on:click={() => push("/pacients-summary")}  class="dh-entry">Pacients</div>
            <div class="dh-entry">Questionarys</div>
            <div class="dh-entry">Schedule</div>
        </div>
    </div>
</div>

<style>
    #hub-content{
        display: flex;
        height: 90vh;
        align-items: center;
    }

    #doctor-hub-entrys {
        display: flex;
        flex-direction: row;
        align-items: center;
        width: 100%;
        height: 60vh;
    }

    .dh-entry {
        cursor: pointer;
        font-family: var(--title-font);
        display: flex;
        width: calc(100% / 3);
        height: 60vh;
        font-size: 1.5em;
        justify-content: center;
        align-items: center;
        transition: all 0.4s ease-out;
    }

    .dh-entry:nth-child(2) {
        border-left: 1px solid var(--theme-color);
        border-right: 1px solid var(--theme-color);
    }

    .dh-entry:hover {
        background-color: var(--theme-color);
        color: white;
    }
</style>