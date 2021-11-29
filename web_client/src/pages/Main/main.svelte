<script>
    import Header from "../../components/Header.svelte";
    import Button from "../../components/Button.svelte";
    import { push } from "svelte-spa-router";
    import LeftContentBackground from "./resources/main-triangle.svg";

    const gotoLogin = profile => {
        push(`/login/${profile}`);
    };
</script>

<div id="main-page" class="page">
    <Header>
        <div on:click={() => gotoLogin('doctor')} class="login-link">Login</div>
    </Header>
    <div id="main-page-content">
        <div id="mp-left-content">
            <div id="mplc-background">
                {@html LeftContentBackground}
            </div>
            <div id="mplc-content">
                <div class="mplc-content-title">
                     <span class="mplc-ct-blue-letter">R</span>emote  <span class="mplc-ct-blue-letter">P</span>acient   <span class="mplc-ct-blue-letter">M</span>onitoring
                </div>
                <div id="mplc-vysagota-description">
                    <p>
                        Vysagota provides professional feature rich rpm for freelance medical proffesionals of any kind. we provide many features to monitor mentally ill
                    </p>
                </div>
                <div id="mplc-button-container">
                    <Button 
                        label="Get Started" 
                        isSquared={true} 
                        width="40%" padding="2vh .5vw" 
                        font_size="1.1em"
                        onClick={() => gotoLogin('pacient')}
                    />
                </div>
            </div>
        </div>
        <div id="mp-background-container"></div>
    </div>
</div>

<style>
    .login-link {
        cursor: pointer;
        font-size: 1.2em;
    }

    #main-page {
        overflow-y: hidden;
    }

    #main-page-content {
        overflow-y: hidden;
        position: relative;
        background-size: cover;
        width: 100%;
        height: 91.5vh;
    }

    #mp-left-content {
        position: relative;
        width: 50%;
        z-index: 2;
    }

    #mp-background-container {
        position: absolute;
        top: 0;
        left: 20%;
        width: 80%;
        height: 100%;
        background-image: url("/main-background.jpg");
        background-size: cover;

    }

    #mplc-background {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        z-index: 1;
    }

    :global(#mplc-background > svg) {
        max-height: 91.5vh;
        text-align: left;
        padding: 0;
    }

    #mplc-content {
        position: absolute;
        display: flex;
        width: 45%;
        height: 91.5vh;
        flex-direction: column;
        font-size: 1.4em;
        justify-content: flex-start;
        top: 0;
        left: 0;
        margin-left: 5%;
        padding: 27% 0 0 5%;
        z-index: 2;
    }

    #mplc-content .mplc-content-title {
        width: 120%;
        font-size: 3rem;
        color: var(--grey-color);
        margin-bottom: 2rem;
        user-select: none;
    }

    .mplc-ct-blue-letter {
        color: var(--lighttheme-color);
    }

    #mplc-vysagota-description {
        margin-bottom: 7em;
        user-select: none;
    }
    
</style>