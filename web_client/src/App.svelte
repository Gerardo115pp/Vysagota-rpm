<script>
    import Router from 'svelte-spa-router';
    import {routes} from './routes';
    import vysagota_client from './classes/VysagotaClient';

</script>



<div id="app-container">
    <Router {routes} />
</div>


<style global>
    @import url('https://fonts.googleapis.com/css2?family=Courgette&family=Dosis:wght@200;400;500&display=swap');
    

    :root {
        --title-font: 'Courgette', serif;
        --body-font: 'Dosis', monospace;
        --theme-color: #054D7A;
        --dark-color: rgb(139, 174, 209);
        --grey-color: #97ADBA;
        --lighttheme-color: #55BBFB;
        --light-color: #F5F5F5;
        --danger: rgb(216, 54, 54);
        --ready: rgb(0, 255, 128);
    }

    :global(body) {
        font-family: var(--body-font);
        color: var(--theme-color);
        margin: 0;
    }

    /* horizontal separator */
    :global(.hs) {
        width: 98%;
        height: 1px;
        background-color: var(--theme-color);
        margin: 0 auto;
    }

</style>